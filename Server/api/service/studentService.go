package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"

	"github.com/mskKandula/oes/api/model"
	xlsx "github.com/tealeg/xlsx/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// emailPayload is the message body published to the "email" RabbitMQ queue.
type emailPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"` // plaintext — for welcome email only
}

var (
	requiredKeys = []string{
		"Name",
		"Email",
		"Mobile",
		"Password",
	}
)

type studentService struct {
	StudentRepository model.StudentRepository
	Publisher         model.Publisher
	publishMu         sync.Mutex   // guards concurrent Publish calls — *amqp.Channel is not goroutine-safe
	redis             *redis.Client // used for Subscribe on code execution result channel
}

// StudentServiceConfig holds dependencies injected into the student service layer.
type StudentServiceConfig struct {
	StudentRepository model.StudentRepository
	// Publisher is used to publish async email/code jobs to the message queue.
	Publisher model.Publisher
	// Redis is used to subscribe to code execution result notifications.
	Redis *redis.Client
}

func NewStudentService(ssc *StudentServiceConfig) model.StudentService {
	return &studentService{
		StudentRepository: ssc.StudentRepository,
		Publisher:         ssc.Publisher,
		redis:             ssc.Redis,
	}
}

func (ss *studentService) CreateStudents(ctx context.Context, byteArray []byte, clientId string) ([]model.Student, error) {

	result, err := excelToJson(byteArray)
	if err != nil {
		return nil, err
	}

	students := make([]model.Student, len(result))

	g, ctx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, runtime.NumCPU()) // bounded concurrency

	for index, val := range result {

		g.Go(func() error {

			sem <- struct{}{}
			defer func() { <-sem }()

			name := val.Get("Name").String()
			email := val.Get("Email").String()
			mobile := val.Get("Mobile").String()
			password := val.Get("Password").String()

			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				return err
			}

			student := model.Student{Id: index + 1, Name: name, Email: email, Mobile: mobile, Password: string(hash), ClientId: clientId}

			if err = ss.StudentRepository.Create(ctx, &student); err != nil {
				return err
			}

			// Publish an async email job to MQServer instead of sending
			// the email synchronously here.
			payload := emailPayload{
				Name:     name,
				Email:    email,
				Mobile:   mobile,
				Password: password, // plaintext carried only in the transient MQ message
			}
			if msgBody, jsonErr := json.Marshal(payload); jsonErr != nil {
				log.Printf("email job: failed to marshal payload for %s: %v", email, jsonErr)
			} else {
				ss.publishMu.Lock()
				pubErr := ss.Publisher.PublishMessageWithContext(ctx, "email", msgBody)
				ss.publishMu.Unlock()
				if pubErr != nil {
					log.Printf("email job: failed to publish for %s: %v", email, pubErr)
				}
			}

			students[index] = student
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return students, nil
}

func excelToJson(fileBytes []byte) ([]gjson.Result, error) {

	var data []gjson.Result

	xlFile, err := xlsx.OpenBinary(fileBytes)

	if err != nil {
		return data, err
	}

	for _, sheet := range xlFile.Sheets {
		if sheet.MaxRow < 2 {
			break
		}

		for rowIndex := 1; rowIndex < sheet.MaxRow; rowIndex++ {

			row, _ := sheet.Row(rowIndex)

			values := []interface{}{}

			for i := 0; i < len(requiredKeys); i++ {

				values = append(values, strings.TrimSpace(row.GetCell(i).String()))

			}

			arr := prepareResult(requiredKeys, values)

			data = append(data, arr)
		}
	}
	return data, nil
}

func prepareResult(keys []string, vals []interface{}) gjson.Result {
	var data string

	for i, k := range keys {
		data, _ = sjson.Set(data, k, vals[i])
	}

	return gjson.Parse(data)
}

// allowedLanguages defines the set of languages accepted for code execution.
var allowedLanguages = map[string]bool{
	"python": true,
	"go":     true,
	"nodejs": true,
}

// SubmitCode validates and publishes a code execution job to RabbitMQ, then waits
// up to 5 seconds for the code-executor microservice to publish a result via Redis
// Pub/Sub. If the result arrives within 5s, it is returned directly (200 OK fast path).
// If not, only the submissionId is returned and the result is delivered via WebSocket
// Type-6 message by the existing oes-server WebSocket hub (202 Accepted slow path).
func (ss *studentService) SubmitCode(ctx context.Context, req model.CodeSubmitRequest, userId, clientId string) (model.CodeSubmitResponse, error) {

	// Validate language
	if !allowedLanguages[req.Language] {
		return model.CodeSubmitResponse{}, fmt.Errorf("unsupported language: %s — allowed: python, go, nodejs", req.Language)
	}

	// Validate code size
	if len(req.Code) > 64*1024 {
		return model.CodeSubmitResponse{}, fmt.Errorf("code exceeds 64KB limit")
	}

	// Default / clamp timeout
	if req.TimeoutMs <= 0 || req.TimeoutMs > 30000 {
		req.TimeoutMs = 10000
	}

	submissionId := uuid.New().String()

	// Subscribe to the result channel BEFORE publishing.
	// Critical ordering: if code-executor finishes and publishes before we subscribe,
	// we would miss the message and always fall through to the 202 path.
	sub := ss.redis.Subscribe(ctx, "result:"+submissionId)
	defer sub.Close()

	// Build the job envelope consumed by code-executor
	job := model.CodeJob{
		SubmissionId: submissionId,
		Language:     req.Language,
		Code:         req.Code,
		Stdin:        req.Stdin,
		TimeoutMs:    req.TimeoutMs,
		UserId:       userId,
		ClientId:     clientId,
	}

	body, err := json.Marshal(job)
	if err != nil {
		return model.CodeSubmitResponse{}, fmt.Errorf("failed to marshal code job: %w", err)
	}

	// Publish to the per-language RabbitMQ queue.
	// Use the same publishMu pattern as CreateStudents() — amqp.Channel is not goroutine-safe.
	ss.publishMu.Lock()
	pubErr := ss.Publisher.PublishMessageWithContext(ctx, "code.execute."+req.Language, body)
	ss.publishMu.Unlock()
	if pubErr != nil {
		return model.CodeSubmitResponse{}, fmt.Errorf("failed to queue code job: %w", pubErr)
	}

	// Mixed pattern: optimistic sync wait with async fallback.
	select {
	case msg := <-sub.Channel():
		// Fast path — code-executor published result within 5s.
		var resp model.CodeSubmitResponse
		if jsonErr := json.Unmarshal([]byte(msg.Payload), &resp); jsonErr != nil {
			return model.CodeSubmitResponse{}, fmt.Errorf("failed to parse execution result: %w", jsonErr)
		}
		resp.Pending = false
		return resp, nil

	case <-time.After(5 * time.Second):
		// Slow path — job is still running; result arrives via WebSocket Type-6 message.
		return model.CodeSubmitResponse{
			SubmissionId: submissionId,
			Pending:      true,
		}, nil
	}
}

func (ss *studentService) FetchStudents(ctx context.Context, clientId string) ([]model.Student, error) {
	return ss.StudentRepository.ReadAll(ctx, clientId)
}

func (ss *studentService) FetchAndPrepare(ctx context.Context, sheetName, clientId string) (*xlsx.File, error) {
	students, err := ss.StudentRepository.ReadAll(ctx, clientId)
	if err != nil {
		return nil, err
	}

	file, err := prepareExcelFile(sheetName, students)
	if err != nil {
		return nil, err
	}

	return file, err
}

func prepareExcelFile(SheetName string, students []model.Student) (*xlsx.File, error) {
	var file *xlsx.File

	var result []map[string]interface{}

	byteData, err := json.Marshal(students)
	if err != nil {
		return file, err
	}

	err = json.Unmarshal(byteData, &result)
	if err != nil {
		return file, err
	}

	file, err = generateExcel(result, SheetName)
	if err != nil {
		return file, err
	}

	return file, nil
}

func generateExcel(studentListResult []map[string]interface{}, SheetName string) (*xlsx.File, error) {

	var (
		file       *xlsx.File
		sheet      *xlsx.Sheet
		row        *xlsx.Row
		rowHeaders []string
	)

	file = xlsx.NewFile()

	sheet, err := file.AddSheet(SheetName)
	if err != nil {
		return file, err
	}

	row = sheet.AddRow()

	row.SetHeight(15)

	row.Hidden = false

	for key, val := range studentListResult[0] {
		if len(val.(string)) > 0 {
			row.AddCell().SetString(strings.ToUpper(key))
			rowHeaders = append(rowHeaders, key)
		}
	}

	for _, obj := range studentListResult {

		row = sheet.AddRow()

		row.SetHeight(15)

		row.Hidden = false

		for _, key := range rowHeaders {

			val := obj[key]

			row.AddCell().SetString(val.(string))
		}
	}
	return file, nil
}

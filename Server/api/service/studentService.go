package service

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/util/emailConfig"
	xlsx "github.com/tealeg/xlsx/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	requiredKeys = []string{
		"Name",
		"Email",
		"Mobile",
		"Password",
	}
	students []model.Student
)

type studentService struct {
	StudentRepository model.StudentRepository
}

// studentServiceCOnfig will hold repositories that will eventually be injected into this
// this service layer
type StudentServiceConfig struct {
	StudentRepository model.StudentRepository
}

func NewStudentService(ssc *StudentServiceConfig) model.StudentService {
	return &studentService{
		StudentRepository: ssc.StudentRepository,
	}
}

func (ss *studentService) CreateStudents(ctx context.Context, byteArray []byte, clientId string) ([]model.Student, error) {

	result, err := excelToJson(byteArray)
	if err != nil {
		return nil, err
	}

	for index, val := range result {

		name := val.Get("Name").String()
		email := val.Get("Email").String()
		mobile := val.Get("Mobile").String()
		password := val.Get("Password").String()

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashedPassword := string(hash)

		id := index + 1

		student := model.Student{id, name, email, mobile, hashedPassword, clientId}

		if err = ss.StudentRepository.Create(ctx, &student); err != nil {
			return nil, err
		}

		if err = emailConfig.SendEmail(student); err != nil {
			log.Println("Error while sending email", err)
		}

		students = append(students, student)
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

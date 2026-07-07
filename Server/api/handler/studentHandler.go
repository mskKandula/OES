package handler

import (
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mskKandula/oes/api/model"
)

type ProofData struct {
	ClientId    string
	UserId      string
	ExamId      string
	ZipFilePath string
}

var (
	ResultPaths = make(chan ProofData, 200)
)

func (h *Handler) StudentsRegister(c *gin.Context) {
	// file, handler, err := c.Request.FormFile("myFile")
	// defer file.Close()

	clientId := c.GetString("clientId")
	ctx := c.Request.Context()

	file, err := c.FormFile("myFile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	base := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(base))

	if ext != ".xlsx" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported File Format"})
		return
	}

	if file.Size > 10*1024 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "File size is big"})
		return
	}

	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer fileData.Close()

	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	students, err := h.StudentService.CreateStudents(ctx, fileBytes, clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": students})
}

func (h *Handler) GetStudents(c *gin.Context) {
	clientId := c.GetString("clientId")
	ctx := c.Request.Context()

	students, err := h.StudentService.FetchStudents(ctx, clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": students})
}

func (h *Handler) DownloadStudents(c *gin.Context) {
	clientId := c.GetString("clientId")
	ctx := c.Request.Context()

	file, err := h.StudentService.FetchAndPrepare(ctx, "All Students Details", clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ReportName := "All Students Details" + ".xlsx"

	c.Writer.Header().Add("Content-type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+ReportName)
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	file.Write(c.Writer)
}

func (h *Handler) GetQuestions(c *gin.Context) {
	// Get examId from query parameter
	examIdStr := c.Query("examId")
	if examIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "examId is required"})
		return
	}

	examId, err := strconv.ParseInt(examIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid examId"})
		return
	}

	// Retrieve questions from cache
	questions, exists := h.QuestionCache.Get(examId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "questions not found for this exam"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions, "examId": examId})
}

func (h *Handler) UploadExamProof(c *gin.Context) {

	bindFile := struct {
		ExamId  int                   `form:"examId" binding:"required"`
		ZipFile *multipart.FileHeader `form:"zipFile" binding:"required"`
	}{}

	// file, err := c.FormFile("zipFile")

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// Bind file
	if err := c.ShouldBind(&bindFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := bindFile.ZipFile
	base := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(base))

	if ext != ".zip" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported File Format"})
		return
	}

	clientId := c.GetString("clientId")
	userId := c.GetInt("userId")

	examStrId := strconv.Itoa(bindFile.ExamId)
	userStrId := strconv.Itoa(userId)

	dstPath := filepath.Join("../media/examProofs", clientId, examStrId, userStrId, base)

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	ResultPaths <- ProofData{
		ClientId:    clientId,
		UserId:      userStrId,
		ExamId:      examStrId,
		ZipFilePath: dstPath,
	}

	c.JSON(http.StatusOK, gin.H{"fileUploaded": "Success"})
}

// AskQuestion handles student Q&A requests against indexed course material via the RAG pipeline.
// POST /api/r/askQuestion
// Body: { "question": "...", "contextId": "optional-topic" }
func (h *Handler) AskQuestion(c *gin.Context) {
	req := model.AskQuestionRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientId := c.GetString("clientId")
	ctx := c.Request.Context()
	answer, err := h.UserService.AskQuestion(ctx, req.Question, req.ContextId, clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}

// ExecuteCode handles student code execution requests.
// POST /r/executeCode
// Body: { "language": "python"|"go"|"nodejs", "code": "...", "stdin": "", "timeoutMs": 5000 }
//
// Response 200 OK (fast path — result within 5s):
//
//	{ "submissionId": "...", "pending": false, "status": "completed", "stdout": "...",
//	  "stderr": "...", "exitCode": 0, "durationMs": 89 }
//
// Response 202 Accepted (slow path — still running):
//
//	{ "submissionId": "...", "pending": true, "message": "..." }
//	Result is delivered asynchronously via WebSocket Type-6 message.
func (h *Handler) ExecuteCode(c *gin.Context) {
	var req model.CodeSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetString("userId")
	clientId := c.GetString("clientId")
	ctx := c.Request.Context()

	resp, err := h.StudentService.SubmitCode(ctx, req, userId, clientId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if resp.Pending {
		c.JSON(http.StatusAccepted, gin.H{
			"submissionId": resp.SubmissionId,
			"pending":      true,
			"message":      "execution in progress, result will be delivered via WebSocket",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

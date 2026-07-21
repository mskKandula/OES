package handler

import (
	"bufio"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mskKandula/oes/api/model"
)

var (
// BufChan       = make(chan string, 10)
// Removed global fileTextLines to prevent memory leaks
// Removed global examId to prevent race conditions
)

func (h *Handler) SignUp(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Users[user.Email] = user.Password
	ctx := c.Request.Context()
	err := h.UserService.CreateUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Successfully signed up"})
}

func (h *Handler) VideoUpload(c *gin.Context) {
	// file, handler, err := c.Request.FormFile("videoFile")

	file, err := c.FormFile("videoFile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// defer file.Close()

	// paths := strings.Split(file.Filename, ".")

	base := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(base))
	fileName := strings.TrimSuffix(base, ext)

	// checking the File Type, if not mp4 return
	if ext != ".mp4" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported File Format"})
		return
	}

	// checking the File Size, if more than 10mb return
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "File size is big"})
		return
	}

	imageName := fileName + ".png"

	clientId := c.GetString("clientId")

	dstPath := filepath.Join("../media/video", clientId, fileName, base)

	m3u8Path := filepath.Join("/media/video", clientId, fileName, "index.m3u8")
	imagePath := filepath.Join("/media/video", clientId, fileName, imageName)

	// Upload the file to specific dst.
	if err = c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	ctx := c.Request.Context()

	if err = h.UserService.CreateVideoFile(ctx, fileName, m3u8Path, imagePath, clientId, dstPath); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	// err = h.UserService.EncodeVideoFile(dstPath)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"fileUploaded": "Success"})
}

func (h *Handler) QuestionsUpload(c *gin.Context) {

	bindFile := struct {
		ExamName     string                `form:"examName" binding:"required"`
		ExamType     string                `form:"examType" binding:"required"`
		QuestionFile *multipart.FileHeader `form:"questionFile" binding:"required"`
	}{}

	// file, err := c.FormFile("myFile")

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// Bind file
	if err := c.ShouldBind(&bindFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := bindFile.QuestionFile
	base := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(base))

	if ext != ".txt" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported File Format"})
		return
	}

	if file.Size > 10*1024 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "File size is big"})
		return
	}

	// buf := bytes.NewBuffer(nil)
	// if _, err = io.Copy(buf, file); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer fileData.Close()

	fileScanner := bufio.NewScanner(fileData)

	fileScanner.Split(bufio.ScanLines)

	// Use local variable instead of global to avoid memory leaks and race conditions
	fileTextLines := make([]string, 0, 100) // Pre-allocate with capacity for better performance
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	clientId := c.GetString("clientId")
	ctx := c.Request.Context()

	examId, err := h.UserService.CreateExam(ctx, clientId, bindFile.ExamName, bindFile.ExamType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to insert exam data"})
		return
	}

	// Cache questions for retrieval by students
	h.QuestionCache.Set(examId, fileTextLines)

	c.JSON(http.StatusOK, gin.H{"questions": fileTextLines, "examId": examId})
}

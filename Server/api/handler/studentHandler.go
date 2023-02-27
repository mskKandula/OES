package handler

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ResultPaths = make(chan string, 200)
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

	if strings.Split(file.Filename, ".")[1] != "xlsx" {
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Add("Content-type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+ReportName+".xlsx")
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	file.Write(c.Writer)
}

func (h *Handler) GetQuestions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"questions": fileTextLines})
}

func (h *Handler) UploadExamProof(c *gin.Context) {
	file, err := c.FormFile("zipFile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.Split(file.Filename, ".")[1] != "zip" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported File Format"})
		return
	}

	dstPath := filepath.Join("../media/examProofs", file.Filename)

	// Upload the file to specific dst.
	// if err = c.SaveUploadedFile(file, dstPath); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
	// 	return
	// }

	// FilePath Creation
	dstFile, err := Create(dstPath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer fileData.Close()

	_, err = io.Copy(dstFile, fileData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer dstFile.Close()

	ResultPaths <- dstPath

	c.JSON(http.StatusOK, gin.H{"fileUploaded": "Success"})
}

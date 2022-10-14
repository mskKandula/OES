package handler

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mskKandula/oes/api/model"
)

var (
	fileTextLines []string
	BufChan       = make(chan string, 10)
)

func (h *Handler) SignUp(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Users[user.Email] = user.Password

	err := h.UserService.CreateUser(user)
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

	paths := strings.Split(file.Filename, ".")

	// checking the File Type, if not mp4 return
	if paths[1] != "mp4" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported File Format"})
		return
	}

	// checking the File Size, if more than 10mb return
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "File size is big"})
		return
	}

	imageName := paths[0] + ".png"

	dstPath := filepath.Join("../media/video", paths[0], file.Filename)

	m3u8Path := filepath.Join("media/video", paths[0], "index.m3u8")
	imagePath := filepath.Join("media/video", paths[0], imageName)

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

	// Upload the file to specific dst.
	// if err = c.SaveUploadedFile(file, dstPath); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
	// 	return
	// }

	if err = h.UserService.CreateVideoFile(paths[0], m3u8Path, imagePath); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	BufChan <- dstPath

	c.JSON(http.StatusOK, gin.H{"fileUploaded": "Success"})
}

func (h *Handler) QuestionsUpload(c *gin.Context) {
	file, err := c.FormFile("myFile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.Split(file.Filename, ".")[1] != "txt" {
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

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	c.JSON(http.StatusOK, gin.H{"questions": fileTextLines})
}

// file path creation
func Create(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return nil, err
	}
	return os.Create(path)
}

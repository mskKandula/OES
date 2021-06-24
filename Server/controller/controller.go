package controller

import (
	"fmt"
	"net/http"
	"bytes"
	"bufio"
	"strings"
	"io"
	"github.com/mskKandula/model"
	"github.com/mskKandula/middleware"
	"github.com/gin-gonic/gin"
)

var Users = make(map[string]string) //temp db

var fileTextLines []string

func SignUp(c *gin.Context){
	user := model.User{}
	if err := c.ShouldBindJSON(&user);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	Users[user.Email]=user.Password
	c.JSON(http.StatusOK,gin.H{"status":"Successfully signed up"})
}

func Login(c *gin.Context){
	userLogin := model.UserLogin{}
	if err := c.ShouldBindJSON(&userLogin);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	if userLogin.Password != Users[userLogin.Email]{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
	}

	token,time,err := middleware.Auth(userLogin)

	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusOK,gin.H{"token":token,"expirationTime":time})
}

func Questionhandle(c *gin.Context){

	file,handler, err := c.Request.FormFile("myFile")
	
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	if strings.Split(handler.Filename, ".")[1] != "txt" {
		fmt.Println("File Format Not supported")
		return
	}

	if handler.Size > 10*1024 {
		fmt.Println("File size is big")
		return
	}

	buf := bytes.NewBuffer(nil)
    if _,err = io.Copy(buf, file); err != nil {
    fmt.Println(err)
	return
}
	fileScanner := bufio.NewScanner(buf)

	fileScanner.Split(bufio.ScanLines)
 
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	c.JSON(http.StatusOK,gin.H{"Questions":fileTextLines})

}

func GetQuestions(w http.ResponseWriter,r *http.Request){
	c.JSON(http.StatusOK,gin.H{"Questions":fileTextLines})
}
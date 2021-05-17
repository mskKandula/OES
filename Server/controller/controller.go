package controller

import (
	"net/http"
	"github.com/mskKandula/model"
	"github.com/mskKandula/middleware"
	"github.com/gin-gonic/gin"
)

var Users = make(map[string]string) //temp db

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
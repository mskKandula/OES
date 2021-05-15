package controller

import (
	"net/http"
	"github.com/mskKandula/model"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context){
	user := model.User{}
	if err := c.ShouldBindJSON(&user);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	c.JSON(http.StatusOK,gin.H{"status":"Successfully signed up"})
}
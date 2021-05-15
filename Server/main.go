package main

import(
	"fmt"
	"github.com/mskKandula/controller"
	"github.com/gin-gonic/gin"
)

func main(){
	fmt.Println("Lets start OES")
	r := gin.Default()
	r.POST("/signUp",controller.SignUp)
	r.Run(":8080")
}
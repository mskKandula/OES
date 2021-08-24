package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mskKandula/controller"
)

func init() {
	controller.Db, err := sql.Open("mysql", "userName:password@tcp(address:port)/OES")

	if err != nil {
		log.Fatalf("Connection Failed to Open: %v",err.Error())
	}

}

func main(){
	fmt.Println("Lets start OES")

	// Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    f, _ := os.Create("Logs/gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

	fs := http.FileServer(http.Dir("../Client/oes/dist"))
	http.Handle("/", fs)

	r := gin.Default()
	r.POST("/signUp",controller.SignUp)
	r.POST("/login",controller.Login)
	r.POST("/multipleStudentsRegistration", controller.StudentsRegisterHandler)
	r.POST("/uploadQuestionFile", controller.QuestionsUploadHandler)
	r.GET("/getRoutes", controller.GetAllRoutes)
	r.GET("/getQuestions",controller.GetQuestions)
	r.GET("/getStudents", controller.GetStudents)
	r.GET("/logOut", controller.Logout)
	r.Run(":8082")
}
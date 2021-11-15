package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mskKandula/controller"
)

var (
	err error
)

func init() {
	controller.Db, err = sql.Open("mysql", "UserName:Password@tcp(127.0.0.1:3306)/OES")

	if err != nil {
		log.Fatalf("Connection Failed to Open: %v", err.Error())
	}

}

func main() {
	fmt.Println("Lets start OES")

	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("Logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// fs := http.FileServer(http.Dir("../Client/oes/dist"))

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("../Client/oes/dist", false)))
	r.POST("/signUp", controller.SignUp)
	r.POST("/login", controller.Login)
	r.POST("/multipleStudentsRegistration", controller.StudentsRegisterHandler)
	r.POST("/uploadQuestionFile", controller.QuestionsUploadHandler)
	r.GET("/ws", controller.Notification)
	r.GET("/getRoutes", controller.GetAllRoutes)
	r.GET("/getQuestions", controller.GetQuestions)
	r.GET("/getStudents", controller.GetStudents)
	r.GET("/logOut", controller.Logout)
	r.Run(":8080")
}

package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mskKandula/config"
	"github.com/mskKandula/controller"
	"github.com/mskKandula/runningProcess"
	"github.com/mskKandula/websock"
)

var (
	err error
)

func init() {
	controller.Db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/OES")

	if err != nil {
		log.Fatalf("Connection Failed to Open: %v", err.Error())
	}

	// Create Redis Client
	config.CreateRedisClient()

}

func main() {
	fmt.Println("Lets start OES")

	go runningProcess.HlsVideoConversion(controller.BufChan)

	pool := websock.NewPool()

	go pool.Start()

	defer func() {
		close(controller.BufChan)
	}()

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
	r.POST("/uploadVideoContent", controller.VideoUploadHandler)
	r.GET("/ws", func(w http.ResponseWriter, r *http.Request) {
		controller.ServeWs(pool, w, r)
	})
	r.GET("/getRoutes", controller.GetAllRoutes)
	r.GET("/getQuestions", controller.GetQuestions)
	r.GET("/getVideos", controller.GetVideos)
	r.GET("/getStudents", controller.GetStudents)
	r.GET("/downloadStudents", controller.DownloadStudents)
	r.GET("/logOut", controller.Logout)
	// go func() {
	// 	r.Run(":8081")
	// }()
	r.Run(":8080")
}

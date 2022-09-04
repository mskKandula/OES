package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mskKandula/config"
	"github.com/mskKandula/controller"
	"github.com/mskKandula/middleware"
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

	// sudo service redis-server start

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

	router := initRouter(pool)
	// go func() {
	// 	r.Run(":8081")
	// }()
	// router.Run(":9000")

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func initRouter(pool *websock.Pool) *gin.Engine {
	r := gin.Default()
	// r.Use(static.Serve("/", static.LocalFile("../Client/oes/dist", false)))
	r.GET("/ws", func(c *gin.Context) {
		controller.ServeWs(pool, c.Writer, c.Request)
	})
	open := r.Group("/o")
	{
		open.POST("/signUp", controller.SignUp)
		open.POST("/login", controller.Login)
	}

	restricted := r.Group("/r").Use(middleware.Auth())
	{
		restricted.POST("/multipleStudentsRegistration", controller.StudentsRegisterHandler)
		restricted.POST("/uploadQuestionFile", controller.QuestionsUploadHandler)
		restricted.POST("/uploadVideoContent", controller.VideoUploadHandler)
		restricted.POST("/uploadExamProof", controller.ExamProofHandler)

		restricted.GET("/getRoutes", controller.GetAllRoutes)
		restricted.GET("/getQuestions", controller.GetQuestions)
		restricted.GET("/getVideos", controller.GetVideos)
		restricted.GET("/getStudents", controller.GetStudents)
		restricted.GET("/downloadStudents", controller.DownloadStudents)
		restricted.GET("/logOut", controller.Logout)
	}

	return r
}

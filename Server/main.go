package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	cnf "github.com/mskKandula/oes/api/config"

	"github.com/mskKandula/oes/api"
	"github.com/mskKandula/oes/api/handler"
)

var (
	err error
)

// func init() {

// 	controller.Db, err = sql.Open("mysql", "root:root@123456789@tcp(127.0.0.1:3306)/OES")

// 	if err != nil {
// 		log.Fatalf("Connection Failed to Open: %v", err.Error())
// 	}

// 	// sudo service redis-server start

// 	// Create Redis Client
// 	config.CreateRedisClient()

// }

func main() {
	fmt.Println("Lets start OES")

	// Increase resources limitations
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	if err = cnf.Setup("config.json"); err != nil {
		log.Fatalf("Setup Failed:%v", err.Error())
	}

	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("Logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// fs := http.FileServer(http.Dir("../Client/oes/dist"))

	defer func() {
		close(handler.BufChan)
		close(handler.ResultPaths)
	}()

	router := api.InitRouter()
	// go func() {
	// 	r.Run(":8081")
	// }()
	// router.Run(":9000")

	s := &http.Server{
		Addr:           ":9000",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

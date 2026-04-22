package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
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

	// Set GOMAXPROCS to utilize all CPU cores for better performance
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Printf("Using %d CPU cores", runtime.NumCPU())

	if err = cnf.Setup("config.json"); err != nil {
		log.Fatalf("Setup Failed:%v", err.Error())
	}

	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("Logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	defer func() {
		close(handler.ResultPaths)
		if f != nil {
			f.Close()
		}
	}()

	router := api.InitRouter()

	// Optimized HTTP server configuration
	s := &http.Server{
		Addr:              ":9000",
		Handler:           router,
		ReadTimeout:       30 * time.Second,  // Reduced from 60s for faster timeout
		WriteTimeout:      30 * time.Second,  // Reduced from 60s for faster timeout
		IdleTimeout:       120 * time.Second, // Added idle timeout
		ReadHeaderTimeout: 10 * time.Second,  // Added header read timeout
		MaxHeaderBytes:    1 << 20,           // 1MB
	}

	log.Printf("Server starting on %s", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

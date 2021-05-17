package main

import(
	"fmt"
	"os"
	"io"
	"github.com/mskKandula/controller"
	"github.com/gin-gonic/gin"
)

func main(){
	fmt.Println("Lets start OES")

	// Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    f, _ := os.Create("Logs/gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.POST("/signUp",controller.SignUp)
	r.POST("/login",controller.Login)
	r.Run(":8082")
}
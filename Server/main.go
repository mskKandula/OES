package main

import(
	"fmt"
	"os"
	"io"
	"github.com/mskKandula/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	controller.Db, err := sql.Open("mysql", "userName:password@tcp(address:port)/oes")

	if err != nil {
		log.Println("Connection Failed to Open")
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
	r.POST("/uploadFile", controller.FileUpload)
	r.POST("/uploadQuestionFile", controller.Questionhandle)
	r.GET("/getQuestions",controller.GetQuestions)
	r.GET("/getStudents", controller.GetStudents)
	r.GET("/logOut", controller.Logout)
	r.Run(":8082")
}
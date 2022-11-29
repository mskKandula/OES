package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mskKandula/oes/api/handler"
	"github.com/mskKandula/oes/api/middleware"
	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/api/repository"
	"github.com/mskKandula/oes/api/service"
	ds "github.com/mskKandula/oes/dataSources"
	"github.com/mskKandula/oes/util/runningProcess"
	"github.com/mskKandula/oes/util/websock"
)

func initSources() (*websock.Pool, *handler.Handler) {
	ds, err := ds.InitDS()
	if err != nil {
		log.Fatalf("Connection Failed to Open:%v", err.Error())
	}

	h := handler.NewHandler(getUserService(ds), getStudentService(ds), getCommonService(ds))

	// go runningProcess.HlsVideoConversion(handler.BufChan)

	go runningProcess.UnzipFile(handler.ResultPaths)

	pool := websock.NewPool()

	go pool.Start(ds)

	return pool, h
}

func InitRouter() *gin.Engine {

	pool, h := initSources()

	r := gin.Default()
	// r.Use(static.Serve("/", static.LocalFile("../Client/oes/dist", false)))

	// r.Use(cor.Default())

	r.GET("/ws", func(c *gin.Context) {
		h.ServeWs(pool, c.Writer, c.Request)
	})

	open := r.Group("/o")
	{
		open.POST("/signUp", h.SignUp)
		open.POST("/login", h.Login)
	}

	common := r.Group("/r").Use(middleware.Auth("Common"))
	{
		common.GET("/getRoutes", h.GetAllRoutes)
		common.GET("/getVideos", h.GetAllVideos)
		common.GET("/logOut", h.Logout)
	}

	user := r.Group("/r").Use(middleware.Auth("Examiner"))
	{
		user.POST("/multipleStudentsRegistration", h.StudentsRegister)
		user.POST("/uploadQuestionFile", h.QuestionsUpload)
		user.POST("/uploadVideoContent", h.VideoUpload)

		user.GET("/getStudents", h.GetStudents)
		user.GET("/downloadStudents", h.DownloadStudents)
	}

	student := r.Group("/r").Use(middleware.Auth("Student"))
	{
		student.POST("/uploadExamProof", h.UploadExamProof)

		student.GET("/getQuestions", h.GetQuestions)

	}

	return r
}

func getUserService(ds *ds.DataSources) model.UserService {

	userMySQLRepository := repository.NewUserMySQLRepository(&repository.RepositoryConfig{
		MySQLDB: ds.MySQLDB, RabbitMQ: ds.RabbitMQ, Queue: ds.Queue})

	userService := service.NewUserService(&service.UserServiceConfig{
		UserRepository: userMySQLRepository,
	})

	return userService
}

func getStudentService(ds *ds.DataSources) model.StudentService {

	studentMySQLRepository := repository.NewStudentMySQLRepository(&repository.RepositoryConfig{
		MySQLDB: ds.MySQLDB})

	studentService := service.NewStudentService(&service.StudentServiceConfig{
		StudentRepository: studentMySQLRepository,
	})

	return studentService
}

func getCommonService(ds *ds.DataSources) model.CommonService {

	commonMySQLRepository := repository.NewCommonMySQLRepository(&repository.RepositoryConfig{
		MySQLDB: ds.MySQLDB, Redis: ds.Redis})

	commonService := service.NewCommonService(&service.CommonServiceConfig{
		CommonRepository: commonMySQLRepository,
	})

	return commonService
}

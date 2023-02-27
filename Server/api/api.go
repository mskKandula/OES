package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mskKandula/oes/api/config"
	"github.com/mskKandula/oes/api/handler"
	"github.com/mskKandula/oes/api/middleware"
	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/api/pkg/questgen/pb"
	"github.com/mskKandula/oes/api/repository"
	"github.com/mskKandula/oes/api/service"
	ds "github.com/mskKandula/oes/dataSources"
	"github.com/mskKandula/oes/util/runningProcess"
	"github.com/mskKandula/oes/util/websock"
	"google.golang.org/grpc"
)

const (
	maxWorkers int = 10
)

func initSources() (*websock.Pool, *handler.Handler) {
	ds, err := ds.InitDS()
	if err != nil {
		log.Fatalf("Connection Failed to Open:%v", err.Error())
	}

	client, err := InitGrpcServiceClient()
	if err != nil {
		log.Fatalf("Connection Failed to Open:%v", err.Error())
	}

	h := handler.NewHandler(getUserService(ds, client), getStudentService(ds), getCommonService(ds))

	// go runningProcess.HlsVideoConversion(handler.BufChan)
	// Worker Pool
	for i := 0; i < maxWorkers; i++ {
		go runningProcess.UnzipFile(handler.ResultPaths)
	}

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
		open.POST("/status", h.CheckStatus)
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
		user.POST("/questionGeneration", h.QuestionGen)

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

func getUserService(ds *ds.DataSources, client pb.QuestGenServiceClient) model.UserService {

	userMySQLRepository := repository.NewUserMySQLRepository(&repository.RepositoryConfig{
		MySQLDB: ds.MySQLDB, RabbitMQ: ds.RabbitMQ, Queue: ds.Queue})

	userService := service.NewUserService(&service.UserServiceConfig{
		UserRepository: userMySQLRepository,
		QuestgenClient: client,
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

func InitGrpcServiceClient() (pb.QuestGenServiceClient, error) {
	// using WithInsecure() because no SSL running
	conn, err := grpc.Dial(config.DatabaseConfig.GRPCDSN, grpc.WithInsecure())
	// defer conn.Close()
	if err != nil {
		return nil, err
	}

	return pb.NewQuestGenServiceClient(conn), nil
}

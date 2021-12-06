package controller

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"
	"github.com/mskKandula/emailConfig"
	"github.com/mskKandula/middleware"
	"github.com/mskKandula/model"
	"github.com/mskKandula/websock"
	"github.com/tealeg/xlsx/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Users = make(map[string]string) //temp db
	Db            *sql.DB
	err           error
	fileTextLines []string
	students      []model.Student
	rowHeaders    []string

	requiredKeys = []string{
		"Name",
		"Email",
		"Mobile",
		"Password",
	}
)

func SignUp(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Users[user.Email] = user.Password

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedPassword := string(hash)

	query, err := Db.Prepare("INSERT INTO Users(name, age, email, mobileNo, password) VALUES(?,?,?,?,?)")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := query.Exec(user.Name, user.Age, user.Email, user.MobileNo, hashedPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lId, _ := result.LastInsertId()

	query, err = Db.Prepare("INSERT INTO UserRole(userId, roleId) VALUES(?,?)")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query.Exec(lId, 1)

	c.JSON(http.StatusOK, gin.H{"status": "Successfully signed up"})
}

func Login(c *gin.Context) {

	userLogin := model.UserLogin{}

	var (
		id       int
		password string
		userType string = "User"
	)

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if userLogin.Password != Users[userLogin.Email] {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }
	row := Db.QueryRow("select id,password from Users where email=?", userLogin.Email)

	err = row.Scan(&id, &password)

	if err != nil {
		if err == sql.ErrNoRows {
			row := Db.QueryRow("select id,password from Students where email=?", userLogin.Email)

			err = row.Scan(&id, &password)

			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
					return
				}
			}
			userType = "Student"
		}
	}

	// if userLogin.Password != password {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userLogin.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenString, expiriesIn, err := middleware.Auth(userLogin, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expiriesIn,
	})

	c.JSON(http.StatusOK, gin.H{"userType": userType})
}

func StudentsRegisterHandler(c *gin.Context) {

	file, handler, err := c.Request.FormFile("myFile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()

	if strings.Split(handler.Filename, ".")[1] != "xlsx" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported File Format"})
		return
	}

	if handler.Size > 10*1024 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "File size is big"})
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := excelToJson(fileBytes)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, val := range result {

		name := val.Get("Name").String()
		email := val.Get("Email").String()
		mobile := val.Get("Mobile").String()
		password := val.Get("Password").String()

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hashedPassword := string(hash)

		query, err := Db.Prepare("INSERT INTO Students(name, email, mobileNo, password) VALUES(?,?,?,?)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result, err := query.Exec(name, email, mobile, hashedPassword)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lId, _ := result.LastInsertId()

		query, err = Db.Prepare("INSERT INTO UserRole(userId, roleId) VALUES(?,?)")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		query.Exec(lId, 2)

		student := model.BasicDetails{
			name,
			email,
			password,
		}

		if err = emailConfig.SendEmail(student); err != nil {
			log.Println("Error while sending email", err)
		}

		students = append(students, model.Student{name, email, mobile, hashedPassword})
	}
	c.JSON(http.StatusOK, gin.H{"students": students})

}

func excelToJson(fileBytes []byte) ([]gjson.Result, error) {

	var data []gjson.Result

	xlFile, err := xlsx.OpenBinary(fileBytes)

	if err != nil {
		return data, err
	}

	for _, sheet := range xlFile.Sheets {
		if sheet.MaxRow < 2 {
			break
		}

		for rowIndex := 1; rowIndex < sheet.MaxRow; rowIndex++ {

			row, _ := sheet.Row(rowIndex)

			allKeys := []string{}

			for _, v := range requiredKeys {
				allKeys = append(allKeys, v)
			}

			values := []interface{}{}

			for i := 0; i < len(allKeys); i++ {

				values = append(values, strings.TrimSpace(row.GetCell(i).String()))

			}

			arr := prepareResult(allKeys, values)

			data = append(data, arr)
		}
	}
	return data, nil
}

func prepareResult(keys []string, vals []interface{}) gjson.Result {
	var data string

	for i, k := range keys {
		data, _ = sjson.Set(data, k, vals[i])
	}

	return gjson.Parse(data)
}

func QuestionsUploadHandler(c *gin.Context) {

	file, handler, err := c.Request.FormFile("myFile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.Split(handler.Filename, ".")[1] != "txt" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported File Format"})
		return
	}

	if handler.Size > 10*1024 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "File size is big"})
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileScanner := bufio.NewScanner(buf)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	c.JSON(http.StatusOK, gin.H{"questions": fileTextLines})

}

func Notification(c *gin.Context) {

	pool := websock.NewPool()

	go pool.Start()

	serveWs(pool, c.Writer, c.Request)

}

func serveWs(pool *websock.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket Endpoint Hit")

	var details websock.Details

	decoder := schema.NewDecoder()

	decoder.Decode(&details, r.URL.Query())
	// if err != nil {
	//     log.Fprintf(w, "%+v\n", err)
	// }

	conn, err := websock.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websock.Client{
		Conn:    conn,
		Pool:    pool,
		Details: &details,
	}

	pool.Register <- client
	client.Read()
}

func GetQuestions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"questions": fileTextLines})
}

func GetStudents(c *gin.Context) {
	rows, err := Db.Query(`SELECT name,email,mobileNo from Students`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()
	students = nil
	for rows.Next() {
		var student model.Student

		if err := rows.Scan(&student.Name, &student.Email, &student.Mobile); err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		students = append(students, student)

	}
	c.JSON(http.StatusOK, gin.H{"students": students})
}

func Logout(c *gin.Context) {

	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
	})

}

func GetAllRoutes(c *gin.Context) {

	cookie, err := c.Request.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		// For any other type of error, return a bad request status
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// Get the JWT string from the cookie
	tokenString := cookie.Value

	id, err := middleware.Decode(tokenString)
	intId := int(id.(float64))
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var routes []model.Route

	// rows, err := Db.Query(`select * from menu where id in(
	// 	select menuId from roleMenu where roleId =(
	// 	select roleId from userRole where userId=(
	// 	select id from examiner where email=?
	// 	)))`, email)

	// if email != "admin@example.org" {
	// 	val = 2
	// } else {
	// 	val = 1
	// }

	// rows, err := Db.Query(`select * from menu where id in(
	// 	select menuId from roleMenu where roleId =(
	// 	select roleId from userRole where userId=?))`, val)

	rows, err := Db.Query(`SELECT m.id,m.name,m.url,m.description FROM UserRole ur
	INNER JOIN RoleMenu rm ON ur.roleId = rm.roleId
	INNER JOIN menu m ON rm.menuId = m.id
	where ur.userId=?`, intId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var route model.Route

		if err := rows.Scan(&route.Id, &route.Name, &route.Url, &route.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		routes = append(routes, route)
	}
	c.JSON(http.StatusOK, gin.H{"routes": routes})
}

func DownloadStudents(c *gin.Context) {

	file, err := PrepareExcelFile("All Students Details")
	ReportName := "All Students Details" + ".xlsx"
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Add("Content-type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+ReportName+".xlsx")
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	file.Write(c.Writer)
}

func PrepareExcelFile(SheetName string) (*xlsx.File, error) {
	var file *xlsx.File

	var result []map[string]interface{}

	byteData, err := json.Marshal(students)

	if err != nil {
		return file, err
	}

	err = json.Unmarshal(byteData, &result)

	if err != nil {
		return file, err
	}

	file, err = generateExcel(result, SheetName)

	if err != nil {
		return file, err
	}

	return file, nil
}

func generateExcel(studentListResult []map[string]interface{}, SheetName string) (*xlsx.File, error) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	file = xlsx.NewFile()

	sheet, err = file.AddSheet(SheetName)

	if err != nil {
		return file, err
	}

	row = sheet.AddRow()

	row.SetHeight(15)

	row.Hidden = false

	for key, val := range studentListResult[0] {
		if len(val.(string)) > 0 {
			row.AddCell().SetString(strings.ToUpper(key))
			rowHeaders = append(rowHeaders, key)
		}
	}

	for _, obj := range studentListResult {

		row = sheet.AddRow()

		row.SetHeight(15)

		row.Hidden = false

		for _, key := range rowHeaders {

			val := obj[key]

			row.AddCell().SetString(val.(string))

		}

	}
	return file, nil
}

package controller

import (
	"bufio"
	"bytes"
	"database/sql"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mskKandula/middleware"
	"github.com/mskKandula/model"
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
	query.Exec(user.Name, user.Age, user.Email, user.MobileNo, hashedPassword)

	c.JSON(http.StatusOK, gin.H{"status": "Successfully signed up"})
}

func Login(c *gin.Context) {

	userLogin := model.UserLogin{}

	var password string

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if userLogin.Password != Users[userLogin.Email] {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }
	row := Db.QueryRow("select password from Users where email=?", userLogin.Email)

	err = row.Scan(&password)

	if err != nil {
		if err == sql.ErrNoRows {
			row := Db.QueryRow("select password from Students where email=?", userLogin.Email)

			err = row.Scan(&password)

			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
					return
				}
			}
		}
	}

	if userLogin.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	tokenString, expiriesIn, err := middleware.Auth(userLogin)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expiriesIn,
	})

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "expirationTime": expiriesIn})
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

	var students []model.Student

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

		query, err := Db.Prepare("INSERT INTO students(name, email, mobileNo, password) VALUES(?,?,?,?)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		query.Exec(name, email, mobile, hashedPassword)

		students = append(students, model.Student{name, email, mobile, password})
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

func GetQuestions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"questions": fileTextLines})
}

func GetStudents(c *gin.Context) {
	var students []model.Student
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

	email, err := middleware.Decode(tokenString)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var routes []model.Route
	var val int

	// rows, err := Db.Query(`select * from menu where id in(
	// 	select menuId from roleMenu where roleId =(
	// 	select roleId from userRole where userId=(
	// 	select id from examiner where email=?
	// 	)))`, email)

	if email != "admin@example.org" {
		val = 2
	} else {
		val = 1
	}

	rows, err := Db.Query(`select * from menu where id in(
		select menuId from roleMenu where roleId =(
		select roleId from userRole where userId=?))`, val)

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

package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mskKandula/oes/api/model"
)

func GenerateJWT(creds model.UserLogin, id int, userType, clientId string) (string, string, time.Time, error) {

	var err error

	//Creating Access Token
	os.Setenv("jwtKey", "7yt65U745TR57lo9h%$fre#$TR43EW") //this should be in an env file

	// Access Token
	atClaims := jwt.MapClaims{}

	atClaims["authorized"] = true

	atClaims["clientId"] = clientId

	atClaims["id"] = id

	atClaims["userType"] = userType

	expirationTime := time.Now().Add(time.Minute * 15)

	atClaims["expireAt"] = expirationTime

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	atokenString, err := atoken.SignedString([]byte(os.Getenv("jwtKey")))

	// Refresh Token
	rtClaims := jwt.MapClaims{}

	rtClaims["id"] = id

	rtClaims["userType"] = userType

	rtClaims["expireAt"] = time.Now().Add(time.Hour * 24)

	rtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	rtokenString, err := rtoken.SignedString([]byte(os.Getenv("jwtKey")))

	if err != nil {
		return "", "", expirationTime, err
	}

	return atokenString, rtokenString, expirationTime, nil
}

func ValidateToken(tokenString, role string) (interface{}, interface{}, interface{}, error) {

	// Initialize a new instance of `Claims`
	claims := jwt.MapClaims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("jwtKey")), nil
	})

	if err != nil {
		return 0, "", "", err
	}

	if !token.Valid {
		return 0, "", "", err
	}

	id := claims["id"]

	userType := claims["userType"]

	clientId := claims["clientId"]

	if userType.(string) != role && role != "Common" {
		return 0, "", "", errors.New("unauthorized to access this resource")
	}

	return id, userType, clientId, nil
}

func Auth(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			// For any other type of error, return a bad request status
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Get the JWT string from the cookie
		tokenString := cookie.Value

		id, userType, clientId, err := ValidateToken(tokenString, role)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		intId := int(id.(float64))
		uType := userType.(string)
		cId := clientId.(string)

		c.Set("userId", intId)
		c.Set("userType", uType)
		c.Set("clientId", cId)

		c.Next()
	}
}

// func CheckUserType(role, userType string) error {
// 	if role != userType {
// 		return errors.New("Unauthorized to access this resource")
// 	}
// 	return nil
// }

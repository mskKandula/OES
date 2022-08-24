package middleware

import (
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mskKandula/model"
)

func GenerateJWT(creds model.UserLogin, id int) (string, time.Time, error) {

	var err error

	//Creating Access Token
	os.Setenv("jwtKey", "7yt65U745TR57lo9h%$fre#$TR43EW") //this should be in an env file

	atClaims := jwt.MapClaims{}

	atClaims["authorized"] = true

	atClaims["email"] = creds.Email

	atClaims["password"] = creds.Password

	atClaims["id"] = id

	expirationTime := time.Now().Add(time.Minute * 5)

	atClaims["expireAt"] = expirationTime

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	tokenString, err := token.SignedString([]byte(os.Getenv("jwtKey")))

	if err != nil {
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil

}

func ValidateToken(tokenString string) (interface{}, error) {

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
		if err == jwt.ErrSignatureInvalid {
			return 0, err
		}
		return 0, err
	}

	if !token.Valid {
		return 0, err
	}

	id := claims["id"]

	return id, nil
}

func Auth() gin.HandlerFunc {
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

			return
		}

		// Get the JWT string from the cookie
		tokenString := cookie.Value

		id, err := ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		intId := int(id.(float64))

		c.Set("userId", intId)

		c.Next()
	}
}

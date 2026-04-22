package middleware

import (
	"errors"
	"net/http"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mskKandula/oes/api/model"
)

var (
	jwtKey     []byte
	jwtKeyOnce sync.Once
)

// Initialize JWT key once at startup
func initJWTKey() {
	jwtKeyOnce.Do(func() {
		// This should ideally come from environment variable or config
		jwtKey = []byte("7yt65U745TR57lo9h%$fre#$TR43EW")
	})
}

func GenerateJWT(creds model.UserLogin, id int, userType, clientId string) (string, time.Time, error) {
	// Initialize JWT key if not already done
	initJWTKey()

	// Create token with pre-populated claims map for better performance
	expirationTime := time.Now().Add(15 * time.Minute)

	atClaims := jwt.MapClaims{
		"authorized": true,
		"clientId":   clientId,
		"id":         id,
		"userType":   userType,
		"exp":        expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil
}

func ValidateToken(tokenString, role string) (interface{}, interface{}, interface{}, error) {
	// Initialize JWT key if not already done
	initJWTKey()

	// Parse token with optimized claims extraction
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, "", "", err
	}

	if !token.Valid {
		return 0, "", "", errors.New("invalid token")
	}

	// Extract claims with existence checks
	id, idOk := claims["id"]
	userType, userTypeOk := claims["userType"]
	clientId, clientIdOk := claims["clientId"]

	if !idOk || !userTypeOk || !clientIdOk {
		return 0, "", "", errors.New("invalid token claims")
	}

	// Check role authorization
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

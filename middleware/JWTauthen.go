package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTauthen() gin.HandlerFunc {
	return func(c *gin.Context) {
		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
		header := c.Request.Header.Get("Authorization")
		tokenString := strings.Replace(header, "Bearer", " ", 1)
		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("token error: %v", token.Header["alg"])
			}
			return hmacSampleSecret, nil
		})
		//check token login
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set example variable
			c.Set("studentId", claims["studentId"])

		} else {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status":  "ineligible",
				"message": err.Error(),
			})
		}

		// before request
		c.Next()
	}
}

package controler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/RangsanB2001/api_rest_golang/connectDB"
	"github.com/RangsanB2001/api_rest_golang/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

type AuthRegister struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	FullName   string `json:"fullname" binding:"required"`
	ImgProfile string `json:"imgprofile" binding:"required"`
}

func RegisStudent(c *gin.Context) {
	var json AuthRegister
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate student
	var existingStudent user.Student
	result := connectDB.Db.Where("username = ? OR full_name = ?", json.Username, json.FullName).First(&existingStudent)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Duplicate student information",
		})
		return
	}

	// Create Students
	encryptPass, err := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to encrypt password",
		})
		return
	}
	newStudent := user.Student{
		Username:   json.Username,
		Password:   string(encryptPass),
		FullName:   json.FullName,
		ImgProfile: json.ImgProfile,
	}
	if err := connectDB.Db.Create(&newStudent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to register student",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Registered successfully",
	})
}

type LoginJson struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for login student
	var checkLogin user.Student
	connectDB.Db.Where("username = ?", json.Username).First(&checkLogin)
	if checkLogin.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Not username information",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(checkLogin.Password), []byte(json.Password))

	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		// token check for 1 minute
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"studentId": checkLogin.ID,
			"exp":       time.Now().Add(time.Minute * 1).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Login Success",
			"token":   tokenString,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Login Failed",
		})
	}
}

package main

import (
	"fmt"

	"github.com/RangsanB2001/api_rest_golang/connectDB"
	"github.com/RangsanB2001/api_rest_golang/controler"
	"github.com/RangsanB2001/api_rest_golang/controler/student"
	"github.com/RangsanB2001/api_rest_golang/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	connectDB.ConnectDB() // Initialize the database connection
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", controler.RegisStudent)
	r.POST("/Login", controler.Login)
	authorized := r.Group("/students", middleware.JWTauthen())
	authorized.GET("/showall", student.ShowAllStudent)
	authorized.GET("/profile", student.Profile)
	r.Run("localhost:8080") // listen and serve on :8080
}

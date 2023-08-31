package student

import (
	"net/http"

	"github.com/RangsanB2001/api_rest_golang/connectDB"
	"github.com/RangsanB2001/api_rest_golang/user"
	"github.com/gin-gonic/gin"
)

func ShowAllStudent(c *gin.Context) {
	var students []user.Student
	connectDB.Db.Find(&students)

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"message":  "show all students success",
		"students": students,
	})
}

func ShowStudent(c *gin.Context) {
	studentID := c.Param("studentId") // รับ studentId จาก URL
	var students []user.Student
	connectDB.Db.Where("id = ?", studentID).Find(&students) // ค้นหา student ที่มี id เท่ากับ studentID

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"message":  "show student success",
		"students": students,
	})
}

func Profile(c *gin.Context) {
	studentId := c.MustGet("studentId").(float64)
	var student []user.Student
	connectDB.Db.First(&student, studentId)
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":   "ok",
		"message":  "show all students success",
		"students": student,
	})
}

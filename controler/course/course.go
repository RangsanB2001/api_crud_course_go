package course

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RangsanB2001/api_rest_golang/connectDB"
	"github.com/RangsanB2001/api_rest_golang/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Coursename string  `json:"coursename"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"imageurl"`
}

func Showcourse(c *gin.Context) {
	var courses []user.Course

	// ดึงข้อมูลคอร์สทั้งหมดที่ไม่ถูกลบ (ไม่ใช้ Unscoped)
	connectDB.Db.Where("deleted_at IS NULL").Find(&courses)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "show course all",
		"courses": courses,
	})
}

func ShowAllcourse(c *gin.Context) {
	var courses []user.Course

	// ดึงข้อมูลคอร์สทั้งหมด (รวมถึงรายการที่ถูกลบด้วย)
	connectDB.Db.Unscoped().Find(&courses) // ใช้ Unscoped() เพื่อให้เรียกข้อมูลที่ถูกลบด้วย

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "show course all",
		"courses": courses,
	})
}

// InsertCourseHandler handles the insertion of a course
func InsertCourseHandler(c *gin.Context) {
	var newCourse user.Course // Use the fully qualified type name

	// Parse JSON request body into the newCourse struct
	if err := c.BindJSON(&newCourse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	// Call the insertCourse function to insert the course
	courseID, err := insertCourse(newCourse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to insert course",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Course inserted successfully",
		"courseid": courseID,
	})
}

// เพิ่มคอร์สใหม่
func insertCourse(newCourse user.Course) (int, error) {
	result := connectDB.Db.Create(&newCourse)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(newCourse.ID), nil
}

// ลบคอร์สโดยใช้ Soft Delete
func deleteCourse(ID int) error {
	var courses []user.Course

	result := connectDB.Db.Where("id = ?", ID).Delete(&courses)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ดึงข้อมูลคอร์สตาม courseid
func getCourseByID(ID int) (*Course, error) {
	var course Course

	result := connectDB.Db.First(&course, ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &course, nil
}

// อัปเดตคอร์ส
func updateCourse(ID int, newCourse Course) (*Course, error) {
	var course Course

	// ค้นหาคอร์สที่ต้องการอัปเดต
	result := connectDB.Db.First(&course, ID)
	if result.Error != nil {
		return nil, result.Error
	}

	// อัปเดตข้อมูลในคอร์ส
	result = connectDB.Db.Model(&course).Updates(newCourse)
	if result.Error != nil {
		return nil, result.Error
	}

	return &course, nil
}

// Delete Course
func DeleteCourseHandler(c *gin.Context) {
	// Get the value of the "courseid" parameter from the URL
	courseIDParam := c.Param("id")

	// Print the parameter value for debugging
	fmt.Println("Course ID Param:", courseIDParam)

	// Convert the course ID parameter to an integer
	ID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		// If the conversion fails, return an error response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid course ID",
		})
		return
	}

	// Call the deleteCourse function to delete the course by ID
	err = deleteCourse(ID)
	if err != nil {
		// If there's an error during deletion, return an error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete course",
		})
		return
	}

	// If everything is successful, return a success response
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Course deleted successfully",
	})
}

// Fetch Course by ID
func GetCourseByIDHandler(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid course ID",
		})
		return
	}

	course, err := getCourseByID(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch course",
		})
		return
	}

	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Course not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"course": course,
	})
}

// Update Course
func UpdateCourseHandler(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("courseid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid course ID",
		})
		return
	}

	var updatedCourse Course
	if err := c.BindJSON(&updatedCourse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	course, err := updateCourse(ID, updatedCourse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update course",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Course updated successfully",
		"course":  course,
	})
}

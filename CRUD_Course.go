package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const coursePath = "courses"
const BasePath = "/api"

type Course struct {
	CourseID     int     `json:"courseid"`
	Coursename   string  `json:"coursename"`
	Price        float64 `json: "price"`
	ImageURL     string  `json: "imageurl"`
	Created_time string  `json: "created_time"`
}

// setup connectdb
func SetupDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/course_it")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

// show course all
func getCourseList() ([]Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := db.QueryContext(ctx, `SELECT 
	courseid,
	coursename,
	price,
	image_url,
	created_time
	FROM course
	`)

	if err != nil {
		log.Panicln(err.Error())
		return nil, err
	}
	defer result.Close()

	courses := make([]Course, 0)
	for result.Next() {
		var course Course
		result.Scan(&course.CourseID, &course.Coursename, &course.Price, &course.ImageURL, &course.Created_time)

		courses = append(courses, course)
	}
	return courses, nil
}

// add course
func insertProduct(course Course) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, `INSERT INTO course (courseid,coursename,price,image_url,created_time) VALUES (?,?,?,?, current_timestamp())
	`, course.CourseID, course.Coursename, course.Price, course.ImageURL)

	if err != nil {
		log.Panicln(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Panicln(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

// delete course
func removeCourse(courseID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, `DELETE FROM course WHERE courseid = ?`, courseID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func updateCourse(courseID int, newCourse Course) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `UPDATE course SET coursename=?, price=?, image_url=? WHERE courseid=?`,
		newCourse.Coursename, newCourse.Price, newCourse.ImageURL, courseID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Now fetch the updated course details
	updatedCourse, err := getCourse(courseID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return updatedCourse, nil
}

// show course
func getCourse(courseid int) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, `SELECT courseid,coursename,price,image_url,created_time FROM course WHERE courseid = ?`, courseid)

	course := &Course{}

	err := row.Scan(
		&course.CourseID,
		&course.Coursename,
		&course.Price,
		&course.ImageURL,
		&course.Created_time,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil

}

// show course all set mathod
func handlerCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		CourseList, err := getCourseList()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(CourseList)

		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)

		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var course Course
		err := json.NewDecoder(r.Body).Decode(&course)

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		courseID, err := insertProduct(course)

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{"courseid":%d}`, courseID)))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handler paht
func handleCourse(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))

	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	courseID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		course, err := getCourse(courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodDelete:
		err := removeCourse(courseID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodPut: // or http.MethodPatch if doing partial updates
		_, err := getCourse(courseID) // Fetch course details by courseID
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Decode the new course data from the request body
		var newCourse Course
		err = json.NewDecoder(r.Body).Decode(&newCourse)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Update the course
		updatedCourse, err := updateCourse(courseID, newCourse)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest) // Use appropriate status code
			return
		}

		// Return the updated course in the response
		j, err := json.Marshal(updatedCourse)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// set path url
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-Custom")
		handler.ServeHTTP(w, r)
	})
}

// setup path Routes
func SetupRoutes(apiBasePath string) {
	// path where get delete id
	courseHandler := http.HandlerFunc(handleCourse)
	coursesHandler := http.HandlerFunc(handlerCourses)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))
}

func main() {
	SetupDB()
	SetupRoutes(BasePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

package connectDB

import (
	"log"
	"os"

	"github.com/RangsanB2001/api_rest_golang/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

func ConnectDB() {
	dsn := os.Getenv("MYSQL_DNS")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("failed to connect database")
	}
	// Migrate the schema
	Db.AutoMigrate(&user.Student{})
	err = Db.AutoMigrate(&user.Course{})
	if err != nil {
		log.Fatal(err)
	}

}

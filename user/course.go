package user

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Coursename string  `json:"coursename"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"imageurl"`
}

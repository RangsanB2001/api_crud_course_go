package user

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Username   string
	Password   string
	FullName   string
	ImgProfile string
}

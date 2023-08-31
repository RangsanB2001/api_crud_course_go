package json

type AuthRegister struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	FullName   string `json:"fullname" binding:"required"`
	ImgProfile string `json:"imgprofile" binding:"required"`
}

type Course struct {
	Coursename string  `json:"coursename"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"imageurl"`
}

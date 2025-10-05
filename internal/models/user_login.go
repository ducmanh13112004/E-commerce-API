package models

type Login struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

package models

type User struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

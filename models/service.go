package models

type Service struct {
	Posts   int `json:"post"`
	Threads int `json:"thread"`
	Users   int `json:"user"`
	Forums  int `json:"forum"`
}

package models

type Forum struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	User  string `json:"user"`
}

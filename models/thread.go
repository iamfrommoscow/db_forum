package models

type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	ID      int    `json:"id"`
	Votes   int    `json:"votes"`
}

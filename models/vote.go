package models

type Vote struct {
	Thread   int    `json:"thread"`
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

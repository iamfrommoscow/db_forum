package models

type Post struct {
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	Message  string `json:"message"`
	Thread   int    `json:"thread"`
	ID       int    `json:"id"`
	Parent   int    `json:"parent"`
	IsEdited bool   `json:"isEdited"`
}

type ReturnPost struct {
	Pst    *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thrd   *Thread `json:"thread,omitempty"`
	Frm    *Forum  `json:"forum,omitempty"`
}

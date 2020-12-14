package models

type Forum struct {
	Posts   int64  `json:"posts"`
	Slug    string `json:"slug" binding:"required"`
	Threads int32  `json:"threads"`
	Title   string `json:"title" binding:"required"`
	User    string `json:"user" binding:"required"`
}


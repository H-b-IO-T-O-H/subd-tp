package models

import (
	"github.com/go-openapi/strfmt"
)

type Post struct {
	ID       int64           `json:"id"`
	Parent   int64           `json:"parent"`
	Author   string          `json:"author"`
	Message  string          `json:"message"`
	IsEdited bool            `json:"isEdited"`
	Forum    string          `json:"forum"`
	Thread   int             `json:"thread"`
	Created  strfmt.DateTime `json:"created"`
}

//easyjson:json
type PostsList []Post

type PostUpdate struct {
	ID      int64  `json:"-"`
	Message string `json:"message"`
}

type PostFull struct {
	Post   *Post   `json:"post,omitempty"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}

type PostNew struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Parent  int64  `json:"parent,omitempty"`
}

type PostGetParams struct {
	PostId     int64
	HaveUser   bool
	HaveForum  bool
	HaveThread bool
}

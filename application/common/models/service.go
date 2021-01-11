package models

type ServiceStatus struct {
	UsersCnt   int64 `json:"user"`
	ForumsCnt  int64 `json:"forum"`
	ThreadsCnt int64 `json:"thread"`
	PostsCnt   int64 `json:"post"`
}

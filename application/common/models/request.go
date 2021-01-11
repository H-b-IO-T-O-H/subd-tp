package models

type QueryParams struct {
	Slug  string
	Limit int
	Since string
	Desc  bool
}

type QueryPostParams struct {
	SlugId string
	Limit  int
	Since  int64
	Sort   string
	Desc   bool
}

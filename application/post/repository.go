package post

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IRepositoryPost interface {
	CreatePost(posts models.PostsList, slugId string) (models.PostsList, errors.Err)
	GetById(id int64) (models.Post, errors.Err)
	GetFull(params models.PostGetParams) (models.PostFull, errors.Err)
	UpdatePost(postUpdate models.PostUpdate) (models.Post, errors.Err)
}

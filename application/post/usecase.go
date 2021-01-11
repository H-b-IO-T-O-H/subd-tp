package post

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IUseCasePost interface {
	CreatePost(posts models.PostsList, slugId string) (models.PostsList, errors.Err)
	GetPost(params models.PostGetParams) (models.PostFull, errors.Err)
	UpdatePost(postUpdate models.PostUpdate) (models.Post, errors.Err)
}

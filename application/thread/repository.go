package forum

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IRepositoryForum interface {
	CreateForum(forumNew models.Forum) errors.Err
	GetBySlug(slug string) (models.Forum, errors.Err)
	GetUsers(uri models.QueryParams) (models.UsersList, errors.Err)
}
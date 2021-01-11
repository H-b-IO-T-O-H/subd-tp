package forum

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IUseCaseForum interface {
	CreateForum(forumNew models.Forum) (models.Forum, errors.Err)
	GetBySlug(slug string) (models.Forum, errors.Err)
	GetUsers(uri models.QueryParams) (models.UsersList, errors.Err)
	//GetThreads(uri models.UriParams) ([]models.User, errors.Err)
}
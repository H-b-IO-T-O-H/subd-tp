package forum

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IUseCaseForum interface {
	CreateForum(forumNew models.Forum) (models.Forum, errors.Err)
	GetBySlug(slug string) (models.Forum, errors.Err)
	GetUsers(query models.QueryParams) (models.UsersList, errors.Err)
	GetThreads(query models.QueryParams) (models.ThreadsList, errors.Err)
}
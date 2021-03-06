package usecase

import (
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/forum"
)

type forumUseCase struct {
	repos forum.IRepositoryForum
}

func NewForumUseCase(repos forum.IRepositoryForum) forum.IUseCaseForum{
	return &forumUseCase{repos: repos}
}

func (f forumUseCase) CreateForum(forumNew models.Forum) (models.Forum, errors.Err) {
	return f.repos.CreateForum(forumNew)
}

func (f forumUseCase) GetBySlug(slug string) (models.Forum, errors.Err) {
	return f.repos.GetBySlug(slug)
}

func (f forumUseCase) GetUsers(query models.QueryParams) (models.UsersList, errors.Err) {
	return f.repos.GetUsers(query)
}

func (f forumUseCase) GetThreads(query models.QueryParams) (models.ThreadsList, errors.Err) {
	return f.repos.GetThreads(query)
}

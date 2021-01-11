package usecase

import (
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/thread"
	"subd/application/thread/delivery"
)

type threadUseCase struct {
	repos thread.IRepositoryThread
}

func NewThreadUseCase(repos thread.IRepositoryThread) thread.IUseCaseThread {
	return &threadUseCase{repos: repos}
}

func (t threadUseCase) CreateThread(thread models.Thread) (models.Thread, errors.Err) {
	id, err := t.repos.CreateThread(thread)
	if err != nil {
		return models.Thread{}, err
	}
	thread.ID = id
	return thread, nil
}

func (t threadUseCase) GetBySlugOrId(slug string, id int) (models.Thread, errors.Err) {
	if id == delivery.BySlug {
		return t.repos.GetBySlug(slug)
	}
	return t.repos.GetById(id)
}


func (t threadUseCase) UpdateBySlugOrId(threadNew models.Thread, slugId string) (models.Thread, errors.Err) {
	return t.repos.UpdateBySlugOrId(threadNew, slugId)
}


func (t threadUseCase) CreateVote(vote models.Vote) (models.Thread, errors.Err) {
	return t.repos.UpsertVote(vote)
}
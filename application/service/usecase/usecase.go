package usecase

import (
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/user"
)

type UserUseCase struct {
	repos user.IRepositoryUser
}

func NewUserUseCase(repos user.IRepositoryUser) user.IUseCaseUser{
	return &UserUseCase{repos: repos}
}

func (u UserUseCase) GetUser(nickname string) (models.User, errors.Err) {
	return u.repos.GetUser(nickname)
}

func (u UserUseCase) UpdateUser(userNew models.User) (models.User, errors.Err) {
	return u.repos.UpdateUser(userNew)
}

func (u UserUseCase) CreateUser(userNew models.User) (models.User, errors.Err) {
	return userNew, u.repos.CreateUser(userNew)
}

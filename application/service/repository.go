package user

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IRepositoryUser interface {
	CreateUser(userNew models.User) errors.Err
	GetUser(nickname string) (models.User, errors.Err)
	UpdateUser(userNew models.User) (models.User, errors.Err)
}
package service

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IRepositoryService interface {
	GetStatus() (models.ServiceStatus, errors.Err)
	Clear() errors.Err
}
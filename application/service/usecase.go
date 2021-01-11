package service

import (
	"subd/application/common/errors"
	"subd/application/common/models"
)

type IUseCaseService interface {
	GetStatus() (models.ServiceStatus, errors.Err)
	Clear() errors.Err
}

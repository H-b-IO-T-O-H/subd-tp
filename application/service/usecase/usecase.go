package usecase

import (
	"subd/application/common/errors"
	"subd/application/common/models"
	"subd/application/service"
)

type ServiceUseCase struct {
	repos service.IRepositoryService
}

func NewServiceUseCase(repos service.IRepositoryService) service.IUseCaseService {
	return &ServiceUseCase{repos: repos}
}

func (s ServiceUseCase) GetStatus() (models.ServiceStatus, errors.Err) {
	return s.repos.GetStatus()
}

func (s ServiceUseCase) Clear() errors.Err {
	return s.repos.Clear()
}

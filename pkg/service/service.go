package service

import (
	"github.com/Dyleme/image-coverter/pkg/repository"
)

type Authorization interface {
}

type History interface {
}

type Service struct {
	repository *repository.Repository
	History
	Authorization
}

func NewService(repository *repository.Repository) *Service {
	return &Service{repository: repository}
}

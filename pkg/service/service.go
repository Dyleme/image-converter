package service

import (
	"github.com/Dyleme/image-coverter"
	"github.com/Dyleme/image-coverter/pkg/repository"
)

type Authorization interface {
	CreateUser(user image.User) (int, error)
	ValidateUser(user image.User) (string, error)
	ParseToken(token string) (int, error)
}

type History interface {
}

type Service struct {
	repository *repository.Repository
	History
	Authorization
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		repository:    rep,
		History:       nil,
		Authorization: NewAuthSevice(*rep),
	}
}

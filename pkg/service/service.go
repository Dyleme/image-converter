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

type Requests interface {
	GetRequests(userID int) ([]image.Request, error)
}

type Service struct {
	repository *repository.Repository
	Authorization
	Requests
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		repository:    rep,
		Requests:      NewRequestService(*rep),
		Authorization: NewAuthSevice(*rep),
	}
}

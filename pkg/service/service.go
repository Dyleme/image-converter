package service

import (
	"mime/multipart"

	"github.com/Dyleme/image-coverter"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/Dyleme/image-coverter/pkg/storage"
)

type Authorization interface {
	CreateUser(user image.User) (int, error)
	ValidateUser(user image.User) (string, error)
	ParseToken(token string) (int, error)
}

type Requests interface {
	GetRequests(userID int) ([]image.Request, error)
	GetRequest(userID int, reqID int) (*image.Request, error)
	AddRequest(int, multipart.File, string, image.ConversionInfo) (int, error)
}

type Service struct {
	repository *repository.Repository
	Authorization
	Requests
	storage.Storage
}

func NewService(rep *repository.Repository, stor *storage.Storage) *Service {
	return &Service{
		repository:    rep,
		Requests:      NewRequestService(*rep, *stor),
		Authorization: NewAuthSevice(*rep),
	}
}

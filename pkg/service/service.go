package service

import (
	"mime/multipart"

	"github.com/Dyleme/image-coverter/pkg/model"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/Dyleme/image-coverter/pkg/storage"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	ValidateUser(user model.User) (string, error)
	ParseToken(token string) (int, error)
}

type Requests interface {
	GetRequests(userID int) ([]model.Request, error)
	GetRequest(userID int, reqID int) (*model.Request, error)
	DeleteRequest(userID int, reqID int) error
	AddRequest(int, multipart.File, string, model.ConversionInfo) (int, error)
}

type Download interface {
	DownloadImage(userID, imageID int) ([]byte, error)
}

type Interface interface {
	Authorization
	Requests
	Download
}

type Service struct {
	Authorization
	Requests
	Download
}

func NewService(rep repository.Interface, stor storage.Interface) *Service {
	return &Service{
		Requests:      NewRequestService(rep, stor),
		Authorization: NewAuthSevice(rep),
		Download:      NewDownloadSerivce(rep, stor),
	}
}

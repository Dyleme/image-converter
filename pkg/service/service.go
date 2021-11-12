package service

import (
	"mime/multipart"

	"github.com/Dyleme/image-coverter/pkg/model"
	"github.com/Dyleme/image-coverter/pkg/repository"
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

type Storager interface {
	GetFile(path string) ([]byte, error)
	UploadFile(userID int, fileName string, data []byte) (string, error)
	DeleteFile(path string) error
}

func NewService(rep repository.Interface, stor Storager) *Service {
	return &Service{
		Requests:      NewRequestService(rep, stor),
		Authorization: NewAuthSevice(rep),
		Download:      NewDownloadSerivce(rep, stor),
	}
}

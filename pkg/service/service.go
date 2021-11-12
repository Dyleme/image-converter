package service

import (
	"context"
	"mime/multipart"

	"github.com/Dyleme/image-coverter/pkg/model"
	"github.com/Dyleme/image-coverter/pkg/repository"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	ValidateUser(ctx context.Context, user model.User) (string, error)
	ParseToken(ctx context.Context, token string) (int, error)
}

type Requests interface {
	GetRequests(ctx context.Context, userID int) ([]model.Request, error)
	GetRequest(ctx context.Context, userID int, reqID int) (*model.Request, error)
	DeleteRequest(ctx context.Context, userID int, reqID int) error
	AddRequest(context.Context, int, multipart.File, string, model.ConversionInfo) (int, error)
}

type Download interface {
	DownloadImage(ctx context.Context, userID, imageID int) ([]byte, error)
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
	GetFile(ctx context.Context, path string) ([]byte, error)
	UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error)
	DeleteFile(ctx context.Context, path string) error
}

func NewService(rep repository.Interface, stor Storager) *Service {
	return &Service{
		Requests:      NewRequestService(rep, stor),
		Authorization: NewAuthSevice(rep),
		Download:      NewDownloadSerivce(rep, stor),
	}
}

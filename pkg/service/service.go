package service

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/Dyleme/image-coverter/pkg/model"
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

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetPasswordAndID(ctx context.Context, nickname string) ([]byte, int, error)
}

type Request interface {
	GetRequests(ctx context.Context, id int) ([]model.Request, error)
	GetRequest(ctx context.Context, userID, reqID int) (*model.Request, error)
	AddRequest(ctx context.Context, req *model.Request, userID int) (int, error)
	DeleteRequest(ctx context.Context, userID, reqID int) (int, int, error)
	AddProcessedImageIDToRequest(ctx context.Context, reqID, imageID int) error
	AddProcessedTimeToRequest(ctx context.Context, reqID int, t time.Time) error
	AddImage(ctx context.Context, userID int, imageInfo model.Info) (int, error)
	DeleteImage(ctx context.Context, userID, imageID int) (string, error)
}

type Download interface {
	GetImageURL(ctx context.Context, userID int, imageID int) (string, error)
}

func NewService(auth, stor Storager) *Service {
	return &Service{
		Requests:      NewRequestService(rep, stor),
		Authorization: NewAuthSevice(rep),
		Download:      NewDownloadSerivce(rep, stor),
	}
}

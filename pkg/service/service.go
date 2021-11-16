package service

import (
	"context"

	"github.com/Dyleme/image-coverter/pkg/handler"
)

type Service struct {
	handler.Autharizater
	handler.Downloader
	handler.Requester
}

type Storager interface {
	GetFile(ctx context.Context, path string) ([]byte, error)
	UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error)
	DeleteFile(ctx context.Context, path string) error
}

func NewService(stor Storager, req Requester, auth Autharizater, down Downloader) *Service {
	return &Service{
		Requester:    NewRequestService(req, stor),
		Autharizater: NewAuthSevice(auth),
		Downloader:   NewDownloadSerivce(down, stor),
	}
}

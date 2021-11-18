package service

import (
	"context"
)

type Storager interface {
	GetFile(ctx context.Context, path string) ([]byte, error)
	UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error)
	DeleteFile(ctx context.Context, path string) error
}

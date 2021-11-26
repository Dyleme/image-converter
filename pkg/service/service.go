package service

import (
	"context"
	"time"
)

type Storager interface {
	GetFile(ctx context.Context, path string) ([]byte, error)
	UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error)
	DeleteFile(ctx context.Context, path string) error
}

type JwtGenerator interface {
	CreateToken(ctx context.Context, tokenTTL time.Duration, id int) (string, error)
}

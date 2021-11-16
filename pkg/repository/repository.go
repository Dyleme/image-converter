package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Dyleme/image-coverter/pkg/model"
)

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

type Repository struct {
	Request
	Authorization
	Download
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Request:       NewReqPostgres(db),
		Download:      NewDownloadPostgres(db),
	}
}

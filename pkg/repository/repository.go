package repository

import (
	"database/sql"
	"time"

	"github.com/Dyleme/image-coverter/pkg/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetPasswordAndID(nickname string) ([]byte, int, error)
}

type Request interface {
	GetRequests(id int) ([]model.Request, error)
	GetRequest(userID, reqID int) (*model.Request, error)
	AddRequest(req *model.Request, userID int) (int, error)
	DeleteRequest(userID, reqID int) (int, int, error)
	AddProcessedImageIDToRequest(reqID, imageID int) error
	AddProcessedTimeToRequest(reqID int, t time.Time) error
	AddImage(userID int, imageInfo model.Info) (int, error)
	DeleteImage(userID, imageID int) (string, error)
}

type Download interface {
	GetImageURL(userID int, imageID int) (string, error)
}

type Interface interface {
	Authorization
	Request
	Download
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

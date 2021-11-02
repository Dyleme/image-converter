package repository

import (
	"database/sql"

	"github.com/Dyleme/image-coverter"
)

type Authorization interface {
	CreateUser(user image.User) (int, error)
	GetPasswordAndID(nickname string) ([]byte, int, error)
}

type Request interface {
	GetRequests(id int) ([]image.Request, error)
	GetRequest(userID, reqID int) (*image.Request, error)
	AddRequest(req *image.Request, userID int) (int, error)
	AddImage(int, image.Info) (int, error)
}

type Repository struct {
	Request
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Request:       NewReqPostgres(db),
	}
}

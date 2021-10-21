package repository

import (
	"database/sql"

	"github.com/Dyleme/image-coverter"
)

type Authorization interface {
	CreateUser(user image.User) (int, error)
	GetPasswordAndID(nickname string) ([]byte, int, error)
}

type History interface {
}

type Repository struct {
	History
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}

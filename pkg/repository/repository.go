package repository

import "database/sql"

type Authorization interface {
}

type History interface {
}

type Repository struct {
	History
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{}
}

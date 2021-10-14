package repository

type Authorization interface {
}

type History interface {
}

type Repository struct {
	History
	Authorization
}

func NewRepository() *Repository {
	return &Repository{}
}

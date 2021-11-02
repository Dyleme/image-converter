package storage

import "os"

type FileStorage interface {
	GetFile(int, string) (*os.File, error)
	UploadFile(int, string, []byte) (string, error)
	DeleteFile(int, string) error
}

type Storage struct {
	FileStorage
}

func NewStorage() *Storage {
	return &Storage{
		FileStorage: NewLocalStorage("D:\\"),
	}
}

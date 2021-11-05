package storage

import "log"

type Interface interface {
	GetFile(path string) ([]byte, error)
	UploadFile(userID int, fileName string, data []byte) (string, error)
	DeleteFile(path string) error
}

type Storage struct {
	Interface
}

func NewStorage() *Storage {
	stor, err := NewMinioStorage("localhost:9000",
		"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", false)
	if err != nil {
		log.Fatalf("can't initialize storage: %v", err)
	}

	return &Storage{Interface: stor}
}

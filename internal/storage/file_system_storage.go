package storage

import (
	"bytes"
	"os"
)

type LocalStorage struct {
	path string
}

func NewLocalStorage(path string) *LocalStorage {
	return &LocalStorage{path: path}
}

func (s *LocalStorage) GetFile(fullPath string) ([]byte, error) {
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	var bf bytes.Buffer

	_, err = bf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return bf.Bytes(), err
}

func (s *LocalStorage) UploadFile(userID int, fileName string, data []byte) (string, error) {
	fullPath := s.path + generateName(fileName)

	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	_, err = file.Write(data)

	return fullPath, err
}

func (s *LocalStorage) DeleteFile(fullPath string) error {
	return os.Remove(fullPath)
}

package storage

import (
	"bytes"
	"os"
)

// LocalSTorage provides methods to store files localy.
type LocalStorage struct {
	path string
}

// NewLoacalStorage is a constructor for LoacalStorage.
func NewLocalStorage(path string) *LocalStorage {
	return &LocalStorage{path: path}
}

// GetFile takes file from the fullPath and returns it's bytes.
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

// UploadFile upload file to the local storage. With the generateated unique name.
// Retuurninig path to this file.
func (s *LocalStorage) UploadFile(userID int, fileName string, data []byte) (string, error) {
	fullPath := s.path + generateName(fileName)

	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	_, err = file.Write(data)

	return fullPath, err
}

// DeleteFile delte file whick path is fullPath.
func (s *LocalStorage) DeleteFile(fullPath string) error {
	return os.Remove(fullPath)
}

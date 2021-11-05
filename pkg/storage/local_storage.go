package storage

import (
	"bytes"
	"os"
	"strconv"
	"strings"
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
	fullPath := s.createPath(userID, fileName)

	for {
		if _, err := os.Stat(fullPath); err != nil {
			break
		}

		var err error

		fullPath, err = s.increaseIndex(fullPath)
		if err != nil {
			return "", err
		}
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	_, err = file.Write(data)

	return fullPath, err
}

func (s *LocalStorage) createPath(userID int, fileName string) string {
	pointIndex := strings.LastIndex(fileName, ".")
	return s.path + strconv.Itoa(userID) + "_" + fileName[:pointIndex] + "_1" + fileName[pointIndex:]
}

func (s *LocalStorage) increaseIndex(path string) (string, error) {
	firstUnder := strings.LastIndex(path, "_")
	secondUnder := strings.LastIndex(path, ".")
	numnber, err := strconv.Atoi(path[firstUnder+1 : secondUnder])

	if err != nil {
		return "", err
	}
	numnber++
	path = path[:firstUnder+1] + strconv.Itoa(numnber) + path[secondUnder:]

	return path, nil
}

func (s *LocalStorage) DeleteFile(fullPath string) error {
	return os.Remove(fullPath)
}

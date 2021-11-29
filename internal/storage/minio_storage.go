package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/minio/minio-go"
)

var ErrBucketNotExist = errors.New("bucket not exist")

type MinioStorage struct {
	client minio.Client
}

func NewMinioStorage(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*MinioStorage, error) {
	cl, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	return &MinioStorage{client: *cl}, nil
}

func (m *MinioStorage) GetFile(ctx context.Context, path string) ([]byte, error) {
	exist, err := m.client.BucketExists("images")

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrBucketNotExist
	}

	obj, err := m.client.GetObject("images", path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	var bf bytes.Buffer
	_, err = bf.ReadFrom(obj)

	if err != nil {
		return nil, err
	}

	return bf.Bytes(), err
}

func (m *MinioStorage) UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error) {
	exist, err := m.client.BucketExists("images")
	if err != nil {
		return "", err
	}

	if !exist {
		e := m.client.MakeBucket("images", "eu-central-1")
		if err != nil {
			return "", e
		}
	}

	var bf bytes.Buffer

	bf.Write(data)

	fileName = m.createPath(userID, fileName)

	for {
		file, _ := m.client.GetObject("images", fileName, minio.GetObjectOptions{})
		if _, err = file.Stat(); err == nil {
			fileName, err = m.increaseIndex(fileName)
			if err != nil {
				return "", fmt.Errorf("minio naming: %w", err)
			}

			continue
		}

		break
	}

	_, err = m.client.PutObject("images", fileName, &bf, int64(bf.Len()), minio.PutObjectOptions{})

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (m *MinioStorage) DeleteFile(ctx context.Context, path string) error {
	exist, err := m.client.BucketExists("images")
	if err != nil {
		return err
	}

	if !exist {
		return ErrBucketNotExist
	}

	err = m.client.RemoveObject("images", path)

	return err
}

func (m *MinioStorage) createPath(userID int, fileName string) string {
	pointIndex := strings.LastIndex(fileName, ".")
	return strconv.Itoa(userID) + "_" + fileName[:pointIndex] + "(1)" + fileName[pointIndex:]
}

func (m *MinioStorage) increaseIndex(path string) (string, error) {
	openBrack := strings.LastIndex(path, "(")
	closeBrack := strings.LastIndex(path, ")")
	numnber, err := strconv.Atoi(path[openBrack+1 : closeBrack])

	if err != nil {
		return "", err
	}
	numnber++
	path = path[:openBrack+1] + strconv.Itoa(numnber) + path[closeBrack:]

	return path, nil
}

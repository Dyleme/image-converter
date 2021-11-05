package storage

import (
	"bytes"
	"errors"

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

func (m *MinioStorage) GetFile(path string) ([]byte, error) {
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

func (m *MinioStorage) UploadFile(userID int, fileName string, data []byte) (string, error) {
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
	_, err = m.client.PutObject("images", fileName, &bf, int64(bf.Len()), minio.PutObjectOptions{})

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (m *MinioStorage) DeleteFile(path string) error {
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

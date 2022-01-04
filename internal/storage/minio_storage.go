package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/minio/minio-go"
)

var ErrBucketNotExist = errors.New("bucket not exist")

// MinioStorage is a struct that provides methods to store files in minio storage.
type MinioStorage struct {
	client minio.Client
}

// MinioConnfig is a config to make connection with minio storage.
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

// NewMinoStorage is a constructor to the MinioStoage.
// Returns error if the connection is denied.
func NewMinioStorage(conf MinioConfig) (*MinioStorage, error) {
	cl, err := minio.New(conf.Endpoint, conf.AccessKeyID, conf.SecretAccessKey, conf.UseSSL)
	if err != nil {
		return nil, fmt.Errorf("can not initialize storage: %w", err)
	}

	return &MinioStorage{client: *cl}, nil
}

// GetFile method get file from minio storage and return it's bytes.
func (m *MinioStorage) GetFile(_ context.Context, path string) ([]byte, error) {
	exist, err := m.client.BucketExists("images")

	if err != nil {
		return nil, fmt.Errorf("cant not get file: %w", err)
	}

	if !exist {
		return nil, ErrBucketNotExist
	}

	obj, err := m.client.GetObject("images", path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("cant not get file: %w", err)
	}

	var bf bytes.Buffer
	_, err = bf.ReadFrom(obj)

	if err != nil {
		return nil, fmt.Errorf("cant not get file: %w", err)
	}

	return bf.Bytes(), err
}

// UploadFile method upload provided file to the minio storage and returns path to the file.
func (m *MinioStorage) UploadFile(_ context.Context, userID int, fileName string, data []byte) (string, error) {
	exist, err := m.client.BucketExists("images")
	if err != nil {
		return "", fmt.Errorf("can not upload file: %w", err)
	}

	if !exist {
		e := m.client.MakeBucket("images", "eu-central-1")
		if err != nil {
			return "", e
		}
	}

	bf := bytes.NewBuffer(data)

	fileName = generateName(fileName)

	_, err = m.client.PutObject("images", fileName, bf, int64(bf.Len()), minio.PutObjectOptions{})

	if err != nil {
		return "", fmt.Errorf("can not upload file: %w", err)
	}

	return fileName, nil
}

// DeleteFile method delete file from the minio storage and return an error if any occurs.
func (m *MinioStorage) DeleteFile(_ context.Context, path string) error {
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

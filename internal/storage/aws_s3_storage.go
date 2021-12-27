package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsStorage struct {
	bucketName string
	session    *session.Session
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

// Initialize AWS S3 storage using environment values.
// Return error if any occurs while initializing session.
func NewAwsStorage(bucketName string, config *aws.Config) (*AwsStorage, error) {
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize session %w", err)
	}

	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	return &AwsStorage{session: sess, uploader: uploader, downloader: downloader, bucketName: bucketName}, nil
}

var ErrReadIsEmpty = errors.New("read empty file")

// GetFile is used to get file from S3 storage.
// Returns an error, any occurs.
func (a *AwsStorage) GetFile(ctx context.Context, path string) ([]byte, error) {
	logger := logging.FromContext(ctx)
	logger.Infof("getting file %v", path)

	downParams := &s3.GetObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(path),
	}

	var b []byte
	buf := aws.NewWriteAtBuffer(b)
	n, err := a.downloader.Download(buf, downParams)

	if err != nil {
		return nil, fmt.Errorf("get file: %w", err)
	}

	if n == 0 {
		return nil, ErrReadIsEmpty
	}

	return buf.Bytes(), nil
}

// UploadFile upload a file to s3 storage.
// Filename is generated like uuid. but file extension ramains the same.
func (a *AwsStorage) UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error) {
	logger := logging.FromContext(ctx)
	logger.Infof("getting file %v", fileName)

	fileName = generateName(fileName)
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewBuffer(data),
	}

	_, err := a.uploader.Upload(upParams)

	if err != nil {
		return "", fmt.Errorf("upload file: %w", err)
	}

	return fileName, nil
}

// DeleteFile delet a file from s3 storage.
// Return an error if any occurs.
func (a *AwsStorage) DeleteFile(ctx context.Context, path string) error {
	svc := s3.New(a.session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &a.bucketName,
		Key:    &path,
	})

	return err
}

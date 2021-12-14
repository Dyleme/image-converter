package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

var bucketName = "dziauho-image-converter"

type AwsStorage struct {
	session    *session.Session
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewAwsStorage() (*AwsStorage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	return &AwsStorage{session: sess, uploader: uploader, downloader: downloader}, nil
}

func (a *AwsStorage) GetFile(ctx context.Context, path string) ([]byte, error) {
	logger := logging.FromContext(ctx)
	logger.Infof("getting file %v", path)

	downParams := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	}

	var b []byte
	buf := aws.NewWriteAtBuffer(b)
	n, err := a.downloader.Download(buf, downParams)

	if n == 0 {
		return nil, fmt.Errorf("didn't read")
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *AwsStorage) UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error) {
	logger := logging.FromContext(ctx)
	logger.Infof("getting file %v", fileName)
	fileName = generateHash(fileName)
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewBuffer(data),
	}

	_, err := a.uploader.Upload(upParams)

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func generateHash(filename string) string {
	dotPos := strings.LastIndex(filename, ".")
	name := uuid.NewString()

	if dotPos != -1 {
		name += filename[dotPos:]
	}

	return name
}

func (a *AwsStorage) DeleteFile(ctx context.Context, path string) error {
	svc := s3.New(a.session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &path,
	})

	return err
}

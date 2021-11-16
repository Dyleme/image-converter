package service

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/Dyleme/image-coverter/pkg/conversion"
	"github.com/Dyleme/image-coverter/pkg/model"

	"image"
)

const (
	jpegType = "jpeg"
	pngType  = "png"
)

const (
	jpegQuality = 100
)

type Requester interface {
	GetRequests(ctx context.Context, id int) ([]model.Request, error)
	GetRequest(ctx context.Context, userID, reqID int) (*model.Request, error)
	AddRequest(ctx context.Context, req *model.Request, userID int) (int, error)
	DeleteRequest(ctx context.Context, userID, reqID int) (int, int, error)
	AddProcessedImageIDToRequest(ctx context.Context, reqID, imageID int) error
	AddProcessedTimeToRequest(ctx context.Context, reqID int, t time.Time) error
	AddImage(ctx context.Context, userID int, imageInfo model.Info) (int, error)
	DeleteImage(ctx context.Context, userID, imageID int) (string, error)
}

type RequestService struct {
	repo    Requester
	storage Storager
}

func NewRequestService(repo Requester, stor Storager) *RequestService {
	return &RequestService{repo: repo, storage: stor}
}

func (s *RequestService) GetRequests(ctx context.Context, userID int) ([]model.Request, error) {
	reqs, err := s.repo.GetRequests(ctx, userID)

	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (s *RequestService) AddRequest(ctx context.Context, userID int, file multipart.File,
	fileName string, convInfo model.ConversionInfo) (int, error) {
	reqTime := time.Now()
	pointIndex := strings.LastIndex(fileName, ".")
	oldType := fileName[pointIndex+1:]

	pic, err := decodeImage(file, oldType)
	if err != nil {
		return 0, err
	}

	url, err := s.UploadImage(ctx, pic, fileName, oldType, userID)
	if err != nil {
		return 0, fmt.Errorf("upload: %w", err)
	}

	x, y := getResolution(pic)
	imageInfo := model.Info{
		ResoultionX: x,
		ResoultionY: y,
		URL:         url,
		Type:        oldType,
	}

	imageID, err := s.repo.AddImage(ctx, userID, imageInfo)
	if err != nil {
		return 0, fmt.Errorf("repo add image: %w", err)
	}

	req := model.Request{
		OpStatus:      "queued",
		RequestTime:   reqTime,
		OriginalID:    imageID,
		Ratio:         convInfo.Ratio,
		OriginalType:  oldType,
		ProcessedType: convInfo.Type,
	}

	reqID, err := s.repo.AddRequest(ctx, &req, userID)
	if err != nil {
		return 0, fmt.Errorf("repo add request: %w", err)
	}

	if convInfo.Ratio != 1 {
		pic = conversion.Convert(pic, convInfo.Ratio)
	}

	convFileName := fileName[:pointIndex] + "_conv." + convInfo.Type

	newURL, err := s.UploadImage(ctx, pic, convFileName, convInfo.Type, userID)
	if err != nil {
		return 0, fmt.Errorf("upload: %w", err)
	}

	newX, newY := getResolution(pic)
	newImageInfo := model.Info{
		ResoultionX: newX,
		ResoultionY: newY,
		URL:         newURL,
		Type:        oldType,
	}

	newImageID, err := s.repo.AddImage(ctx, userID, newImageInfo)
	if err != nil {
		return 0, fmt.Errorf("repo add image: %w", err)
	}

	err = s.repo.AddProcessedImageIDToRequest(ctx, reqID, newImageID)
	if err != nil {
		return 0, fmt.Errorf("repo update image in request: %w", err)
	}

	completionTime := time.Now()

	err = s.repo.AddProcessedTimeToRequest(ctx, reqID, completionTime)
	if err != nil {
		return 0, fmt.Errorf("repo update time in request: %w", err)
	}

	return reqID, nil
}

func decodeImage(r io.Reader, oldType string) (image.Image, error) {
	switch oldType {
	case pngType:
		return png.Decode(r)
	case jpegType:
		return jpeg.Decode(r)
	default:
		return nil, fmt.Errorf("can not work with this type")
	}
}

func getResolution(i image.Image) (x, y int) {
	return i.Bounds().Dx(), i.Bounds().Dy()
}

func encodeImage(i image.Image, fileType string) ([]byte, error) {
	bf := new(bytes.Buffer)

	switch fileType {
	case pngType:
		if err := png.Encode(bf, i); err != nil {
			return nil, err
		}
	case jpegType:
		if err := jpeg.Encode(bf, i, &jpeg.Options{Quality: jpegQuality}); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type of image")
	}

	return bf.Bytes(), nil
}

func (s *RequestService) GetRequest(ctx context.Context, userID, reqID int) (*model.Request, error) {
	req, err := s.repo.GetRequest(ctx, userID, reqID)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *RequestService) DeleteRequest(ctx context.Context, userID, reqID int) error {
	im1ID, im2ID, err := s.repo.DeleteRequest(ctx, userID, reqID)

	if err != nil {
		return err
	}

	url1, err := s.repo.DeleteImage(ctx, userID, im1ID)
	if err != nil {
		return err
	}

	url2, err := s.repo.DeleteImage(ctx, userID, im2ID)
	if err != nil {
		return err
	}

	err = s.storage.DeleteFile(ctx, url1)
	if err != nil {
		return err
	}

	err = s.storage.DeleteFile(ctx, url2)
	if err != nil {
		return err
	}

	return err
}

func (s *RequestService) UploadImage(ctx context.Context, i image.Image,
	fileName, imageType string, userID int) (string, error) {
	bf, err := encodeImage(i, imageType)
	if err != nil {
		return "", err
	}

	newURL, err := s.storage.UploadFile(ctx, userID, fileName, bf)
	if err != nil {
		return "", err
	}

	return newURL, nil
}

package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
	"time"

	"github.com/Dyleme/image-coverter/internal/conversion"
	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"

	"image"
)

const (
	jpegType = "jpeg"
	pngType  = "png"
)

const (
	jpegQuality = 100
)

var ErrUnsupportedType = errors.New("unsopported type")

type Requester interface {
	GetRequests(ctx context.Context, id int) ([]model.Request, error)
	GetRequest(ctx context.Context, userID, reqID int) (*model.Request, error)
	AddRequest(ctx context.Context, req *model.Request, userID int) (int, error)
	DeleteRequest(ctx context.Context, userID, reqID int) (int, int, error)
	UpdateRequestStatus(ctx context.Context, reqID int, status string) error
	AddProcessedImageIDToRequest(ctx context.Context, reqID, imageID int) error
	AddProcessedTimeToRequest(ctx context.Context, reqID int, t time.Time) error
	AddImage(ctx context.Context, userID int, imageInfo model.Info) (int, error)
	DeleteImage(ctx context.Context, userID, imageID int) (string, error)
}

type RequestService struct {
	repo      Requester
	storage   Storager
	processor ImageProcesser
}

type ImageProcesser interface {
	ProcessImage(data *model.ConversionData)
}

func NewRequestService(repo Requester, stor Storager, proc ImageProcesser) *RequestService {
	return &RequestService{repo: repo, storage: stor, processor: proc}
}

func (s *RequestService) GetRequests(ctx context.Context, userID int) ([]model.Request, error) {
	reqs, err := s.repo.GetRequests(ctx, userID)

	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (s *RequestService) AddRequest(ctx context.Context, userID int, file io.Reader,
	fileName string, convInfo model.ConversionInfo) (int, error) {
	reqTime := time.Now()
	pointIndex := strings.LastIndex(fileName, ".")

	if pointIndex == -1 {
		return 0, fmt.Errorf("no point in filename")
	}

	oldType := fileName[pointIndex+1:]

	fileData, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	pic, err := decodeImage(bytes.NewBuffer(fileData), oldType)
	if err != nil {
		return 0, err
	}

	url, err := s.uploadFile(ctx, fileData, fileName, userID)
	if err != nil {
		return 0, fmt.Errorf("upload: %w", err)
	}

	width, height := getResolution(pic)
	imageInfo := model.Info{
		Width:  width,
		Height: height,
		URL:    url,
		Type:   oldType,
	}

	imageID, err := s.repo.AddImage(ctx, userID, imageInfo)
	if err != nil {
		return 0, fmt.Errorf("repo add image: %w", err)
	}

	req := model.Request{
		OpStatus:      repository.StatusQueued,
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

	convertImageData := model.ConversionData{
		Ctx:       ctx,
		ImageInfo: convInfo,
		UserID:    userID,
		ReqID:     reqID,
		OldType:   oldType,
		Pic:       fileData,
		FileName:  fileName,
	}

	s.processor.ProcessImage(&convertImageData)

	return reqID, nil
}

func decodeImage(r io.Reader, oldType string) (image.Image, error) {
	switch oldType {
	case pngType:
		return png.Decode(r)
	case jpegType:
		return jpeg.Decode(r)
	default:
		return nil, ErrUnsupportedType
	}
}

func getResolution(i image.Image) (width, height int) {
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
		return nil, ErrUnsupportedType
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

func (s *RequestService) uploadFile(ctx context.Context, bts []byte,
	fileName string, userID int) (string, error) {
	newURL, err := s.storage.UploadFile(ctx, userID, fileName, bts)
	if err != nil {
		return "", err
	}

	return newURL, nil
}

func (s *RequestService) Convert(ctx context.Context, data *model.ConversionData) image.Image {
	logger := logging.FromContext(ctx)
	err := s.repo.UpdateRequestStatus(ctx, data.ReqID, repository.StatusProcessing)

	if err != nil {
		logger.Warn(fmt.Errorf("repo update status in request: %w", err))
	}

	logger.WithField("name", data.FileName).Info("start image conversion")

	im, err := decodeImage(bytes.NewBuffer(data.Pic), data.OldType)

	if err != nil {
		logger.Error(err)
	}

	if data.ImageInfo.Ratio != 1 {
		im = conversion.Resize(im, data.ImageInfo.Ratio)
	}

	return im
}

func (s *RequestService) ProcessResizedImage(ctx context.Context, im image.Image, data *model.ConversionData) {
	logger := logging.FromContext(ctx)
	pointIndex := strings.LastIndex(data.FileName, ".")
	convFileName := data.FileName[:pointIndex] + "_conv." + data.ImageInfo.Type

	bts, err := encodeImage(im, data.ImageInfo.Type)
	if err != nil {
		logger.Warn(fmt.Errorf("encode image: %w", err))
	}

	newURL, err := s.uploadFile(ctx, bts, convFileName, data.UserID)
	if err != nil {
		logger.Warn(fmt.Errorf("upload: %w", err))
	}

	newX, newY := getResolution(im)
	newImageInfo := model.Info{
		Width:  newX,
		Height: newY,
		URL:    newURL,
		Type:   data.ImageInfo.Type,
	}

	newImageID, err := s.repo.AddImage(ctx, data.UserID, newImageInfo)
	if err != nil {
		logger.Warn(fmt.Errorf("repo add image: %w", err))
	}

	err = s.repo.AddProcessedImageIDToRequest(ctx, data.ReqID, newImageID)
	if err != nil {
		logger.Warn(fmt.Errorf("repo update image in request: %w", err))
	}

	completionTime := time.Now()

	err = s.repo.AddProcessedTimeToRequest(ctx, data.ReqID, completionTime)
	if err != nil {
		logger.Warn(fmt.Errorf("repo update time in request: %w", err))
	}

	err = s.repo.UpdateRequestStatus(ctx, data.ReqID, repository.StatusDone)
	if err != nil {
		logger.Warn(fmt.Errorf("repo update status in request: %w", err))
	}
}

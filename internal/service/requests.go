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
	"github.com/sirupsen/logrus"

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

type ConvesionData struct {
	ctx       context.Context
	imageInfo model.ConversionInfo
	userID    int
	reqID     int
	pic       image.Image
	fileName  string
}

type RequestService struct {
	repo    Requester
	storage Storager
	req     chan *ConvesionData
}

func (s *RequestService) worker(ch <-chan *ConvesionData) {
	for data := range ch {
		s.convert(data)
	}
}

func (s *RequestService) convert(data *ConvesionData) {
	ctx := data.ctx
	logger := logging.FromContext(ctx)

	err := s.repo.UpdateRequestStatus(ctx, data.reqID, repository.StatusProcessing)
	if err != nil {
		logger.Warn(fmt.Errorf("repo update status in request: %w", err))
	}

	begin := time.Now()

	logger.WithField("name", data.fileName).Info("start image conversion")

	if data.imageInfo.Ratio != 1 {
		data.pic = conversion.Convert(data.pic, data.imageInfo.Ratio)
	}

	pointIndex := strings.LastIndex(data.fileName, ".")
	convFileName := data.fileName[:pointIndex] + "_conv." + data.imageInfo.Type

	bts, err := encodeImage(data.pic, data.imageInfo.Type)
	if err != nil {
		logger.Warn(fmt.Errorf("encode image: %w", err))
	}

	newURL, err := s.uploadFile(ctx, bts, convFileName, data.userID)
	if err != nil {
		logger.Warn(fmt.Errorf("upload: %w", err))
	}

	newX, newY := getResolution(data.pic)
	newImageInfo := model.Info{
		ResoultionX: newX,
		ResoultionY: newY,
		URL:         newURL,
		Type:        data.imageInfo.Type,
	}

	newImageID, err := s.repo.AddImage(ctx, data.userID, newImageInfo)
	if err != nil {
		logger.Warn(fmt.Errorf("repo add image: %w", err))
	}

	err = s.repo.AddProcessedImageIDToRequest(ctx, data.reqID, newImageID)
	if err != nil {
		logger.Warn(fmt.Errorf("repo update image in request: %w", err))
	}

	completionTime := time.Now()

	err = s.repo.AddProcessedTimeToRequest(ctx, data.reqID, completionTime)
	if err != nil {
		logger.Warn(fmt.Errorf("repo update time in request: %w", err))
	}

	err = s.repo.UpdateRequestStatus(ctx, data.reqID, repository.StatusDone)
	if err != nil {
		logger.Warn(fmt.Errorf("repo update status in request: %w", err))
	}

	logger.WithFields(logrus.Fields{
		"time for conversion": time.Since(begin),
		"name":                data.fileName,
	}).Info("end image conversion")
}

func NewRequestService(repo Requester, stor Storager, workersAmount uint) *RequestService {
	s := &RequestService{repo: repo, storage: stor, req: make(chan *ConvesionData)}
	for i := 0; i < int(workersAmount); i++ {
		go s.worker(s.req)
	}

	return s
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

	convertImageInfo := ConvesionData{
		ctx:       ctx,
		imageInfo: convInfo,
		userID:    userID,
		reqID:     reqID,
		pic:       pic,
		fileName:  fileName,
	}

	go func() { s.req <- &convertImageInfo }()

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

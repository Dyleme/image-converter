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
	"github.com/Dyleme/image-coverter/internal/rabbitmq"
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

// RequestRepo is an interface which provides methods to implement with the reposistory.
type RequestRepo interface {
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

// Request is a struct provides the abitility to get, add, delete and update requests.
type Request struct {
	repo      RequestRepo
	storage   Storager
	processor ImageProcesser
}

// ImageProcesser is an interface which is provides method to save image to the repo.
type ImageProcesser interface {
	ProcessImage(ctx context.Context, data *rabbitmq.ConversionData)
}

// NewRequest is a constructor to the RequestService.
func NewRequest(repo RequestRepo, stor Storager, proc ImageProcesser) *Request {
	return &Request{repo: repo, storage: stor, processor: proc}
}

// GetRequests returns requsts, or error if any occurs.
// Function get requests with repo.GetRequests and returns them.
func (s *Request) GetRequests(ctx context.Context, userID int) ([]model.Request, error) {
	reqs, err := s.repo.GetRequests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get reqeusts: %w", err)
	}

	return reqs, nil
}

type RatioNotInRangeError struct {
	ratio float32
}

func (e *RatioNotInRangeError) Error() string {
	return fmt.Sprintf("ration should be between 0 and 1, ratio is %v", e.ratio)
}

type FilenameWithoutPotinError struct {
	filename string
}

func (e *FilenameWithoutPotinError) Error() string {
	return fmt.Sprintf("filename should include point, filename is %s", e.filename)
}

// AddRequest return the id of the added request or error if any occurs.
// Also this function calls processor.ProcessImgae to convert the image.
// Function decode file as image and upload this image using stor.UploadFile,
// add request to the repo with repo.AddRequest.
func (s *Request) AddRequest(ctx context.Context, userID int, file io.Reader,
	fileName string, convInfo model.ConversionInfo) (int, error) {
	if convInfo.Ratio > 1 || convInfo.Ratio <= 0 {
		return 0, &RatioNotInRangeError{convInfo.Ratio}
	}

	reqTime := time.Now()

	pointIndex := strings.LastIndex(fileName, ".")
	if pointIndex == -1 {
		return 0, &FilenameWithoutPotinError{fileName}
	}

	oldType := fileName[pointIndex+1:]

	fileData, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	img, err := decodeImage(bytes.NewBuffer(fileData), oldType)
	if err != nil {
		return 0, err
	}

	url, err := s.uploadFile(ctx, fileData, fileName, userID)
	if err != nil {
		return 0, fmt.Errorf("upload: %w", err)
	}

	width, height := getResolution(img)
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

	convertImageData := rabbitmq.ConversionData{
		Ctx:       ctx,
		ImageInfo: convInfo,
		UserID:    userID,
		ReqID:     reqID,
		OldType:   oldType,
		Pic:       fileData,
		FileName:  fileName,
	}

	s.processor.ProcessImage(ctx, &convertImageData)

	return reqID, nil
}

// decodeImage decodes image from the r.
// Decoding supports only jpeg and png types.
func decodeImage(r io.Reader, imgType string) (image.Image, error) {
	switch imgType {
	case pngType:
		return png.Decode(r)
	case jpegType:
		return jpeg.Decode(r)
	default:
		return nil, ErrUnsupportedType
	}
}

// getResolution function returns the resolution of the image.
func getResolution(i image.Image) (width, height int) {
	return i.Bounds().Dx(), i.Bounds().Dy()
}

// encodeImage encode image with the provided image type, returns bytes of the encoded image.
func encodeImage(i image.Image, imgType string) ([]byte, error) {
	bf := new(bytes.Buffer)

	switch imgType {
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

// GetRequest returns the request by its id and user id.
// Method calls repo.GetRequest and return it's result.
func (s *Request) GetRequest(ctx context.Context, userID, reqID int) (*model.Request, error) {
	req, err := s.repo.GetRequest(ctx, userID, reqID)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// DeleteRequest method deletes request.
// At first it deletes request using repo.DeleteRequest, than delete image from database using.DeleteImage
// and finally it deletes images from the storage using storage.DeletFile.
func (s *Request) DeleteRequest(ctx context.Context, userID, reqID int) error {
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

func (s *Request) uploadFile(ctx context.Context, bts []byte,
	fileName string, userID int) (string, error) {
	newURL, err := s.storage.UploadFile(ctx, userID, fileName, bts)
	if err != nil {
		return "", err
	}

	return newURL, nil
}

// Convert is function that converts image, that is getted from ConversionData.
func (s *Request) Convert(ctx context.Context, data *rabbitmq.ConversionData) image.Image {
	logger := logging.FromContext(ctx)

	err := s.repo.UpdateRequestStatus(ctx, data.ReqID, repository.StatusProcessing)
	if err != nil {
		logger.Warn(fmt.Errorf("repo update: %w", err))
	}

	logger.WithField("name", data.FileName).Info("start image conversion")

	begin := time.Now()

	im, err := decodeImage(bytes.NewBuffer(data.Pic), data.OldType)
	if err != nil {
		logger.Error(err)
	}

	if data.ImageInfo.Ratio != 1 {
		im = conversion.Resize(im, data.ImageInfo.Ratio)
	}

	logger.WithFields(logrus.Fields{
		"name":          data.FileName,
		"time for conv": time.Since(begin),
	}).Info("end image conversion")

	return im
}

// ProcessResizedImage is used to upload image to the storage and update repository.
func (s *Request) ProcessResizedImage(ctx context.Context, im image.Image, data *rabbitmq.ConversionData) {
	logger := logging.FromContext(ctx)
	pointIndex := strings.LastIndex(data.FileName, ".")
	convFileName := data.FileName[:pointIndex] + "_conv." + data.ImageInfo.Type

	bts, err := encodeImage(im, data.ImageInfo.Type)
	if err != nil {
		logger.Errorf("encode image: %s", err)
	}

	newURL, err := s.uploadFile(ctx, bts, convFileName, data.UserID)
	if err != nil {
		logger.Errorf("upload: %s", err)
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
		logger.Errorf("repo add image: %s", err)
	}

	err = s.repo.AddProcessedImageIDToRequest(ctx, data.ReqID, newImageID)
	if err != nil {
		logger.Errorf("repo update image in request: %s", err)
	}

	completionTime := time.Now()

	err = s.repo.AddProcessedTimeToRequest(ctx, data.ReqID, completionTime)
	if err != nil {
		logger.Errorf("repo update time in request: %s", err)
	}

	err = s.repo.UpdateRequestStatus(ctx, data.ReqID, repository.StatusDone)
	if err != nil {
		logger.Errorf("repo update status in request: %s", err)
	}
}

package service

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/Dyleme/image-coverter"
	"github.com/Dyleme/image-coverter/pkg/conversion"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/Dyleme/image-coverter/pkg/storage"

	im "image"
)

const (
	jpegType = "jpeg"
	pngType  = "png"
)

const (
	jpegQuality = 100
)

type RequestService struct {
	repo    repository.Request
	storage storage.Storage
}

func NewRequestService(repo repository.Request, stor storage.Storage) *RequestService {
	return &RequestService{repo: repo, storage: stor}
}

func (s *RequestService) GetRequests(userID int) ([]image.Request, error) {
	reqs, err := s.repo.GetRequests(userID)

	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (s *RequestService) AddRequest(userID int, file multipart.File,
	fileName string, info image.ConversionInfo) (int, error) {
	reqTime := time.Now()
	pointIndex := strings.LastIndex(fileName, ".")
	oldType := fileName[pointIndex+1:]

	pic, err := decodeImage(file, oldType)
	if err != nil {
		return 0, err
	}

	url, err := s.UploadImage(pic, fileName, oldType, userID)
	if err != nil {
		return 0, err
	}

	x, y := getResolution(pic)
	imageInfo := image.Info{
		ResoultionX: x,
		ResoultionY: y,
		URL:         url,
		Type:        oldType,
	}

	imageID, err := s.repo.AddImage(userID, imageInfo)
	if err != nil {
		return 0, err
	}

	req := image.Request{
		OpStatus:      "queued",
		RequestTime:   reqTime,
		OriginalID:    imageID,
		Ratio:         info.Ratio,
		OriginalType:  oldType,
		ProcessedType: info.Type,
	}

	_, err = s.repo.AddRequest(&req, userID)
	if err != nil {
		return 0, err
	}

	if info.Ratio != 1 {
		pic = conversion.Convert(pic, info.Ratio)
	}

	convFileName := fileName[:pointIndex] + "_conv." + info.Type

	newURL, err := s.UploadImage(pic, convFileName, info.Type, userID)
	if err != nil {
		return 0, err
	}

	newX, newY := getResolution(pic)
	newImageInfo := image.Info{
		ResoultionX: newX,
		ResoultionY: newY,
		URL:         newURL,
		Type:        oldType,
	}

	newImageID, err := s.repo.AddImage(userID, newImageInfo)
	if err != nil {
		return 0, err
	}

	return newImageID, nil
}

func decodeImage(r io.Reader, oldType string) (im.Image, error) {
	switch oldType {
	case pngType:
		return png.Decode(r)
	case jpegType:
		return jpeg.Decode(r)
	default:
		return nil, fmt.Errorf("can not work with this type")
	}
}

func getResolution(i im.Image) (x, y int) {
	return i.Bounds().Dx(), i.Bounds().Dy()
}

func encodeImage(i im.Image, fileType string) ([]byte, error) {
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

func (s *RequestService) GetRequest(userID, reqID int) (*image.Request, error) {
	req, err := s.repo.GetRequest(userID, reqID)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *RequestService) UploadImage(i im.Image, fileName, imageType string, userID int) (string, error) {
	bf, err := encodeImage(i, imageType)
	if err != nil {
		return "", err
	}

	newURL, err := s.storage.UploadFile(userID, fileName, bf)
	if err != nil {
		return "", err
	}

	return newURL, nil
}

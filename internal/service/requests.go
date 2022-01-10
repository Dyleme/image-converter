package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"
)

// RequestRepo is an interface which provides methods to implement with the reposistory.
type RequestRepo interface {
	GetRequests(ctx context.Context, id int) ([]model.Request, error)
	GetRequest(ctx context.Context, userID, reqID int) (*model.Request, error)
	AddImageAndRequest(ctx context.Context, userID int, imageInfo *model.ReuquestImageInfo,
		req *model.Request) (int, error)
	DeleteRequestAndImage(ctx context.Context, userID, reqID int) (im1url, im2url string, err error)
}

// Request is a struct provides the abitility to get, add, delete and update requests.
type Request struct {
	repo      RequestRepo
	storage   Storager
	processor ImageProcesser
}

// ImageProcesser is an interface which is provides method to save image to the repo.
type ImageProcesser interface {
	ProcessImage(ctx context.Context, data *model.RequestToProcess) error
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

type FilenameWithoutPotintError struct {
	filename string
}

func (e *FilenameWithoutPotintError) Error() string {
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
		return 0, &FilenameWithoutPotintError{fileName}
	}

	oldType := fileName[pointIndex+1:]
	if oldType != jpegType && oldType != pngType {
		return 0, fmt.Errorf("add request: %w", UnsupportedTypeError{oldType})
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	url, err := s.uploadFile(ctx, fileData, fileName, userID)
	if err != nil {
		return 0, fmt.Errorf("add request: %w", err)
	}

	imageInfo := model.ReuquestImageInfo{
		URL:  url,
		Type: oldType,
	}

	req := model.Request{
		OpStatus:      repository.StatusQueued,
		RequestTime:   reqTime,
		Ratio:         convInfo.Ratio,
		OriginalType:  oldType,
		ProcessedType: convInfo.Type,
	}

	reqID, err := s.repo.AddImageAndRequest(ctx, userID, &imageInfo, &req)
	if err != nil {
		return 0, fmt.Errorf("repo add image and request: %w", err)
	}

	convertImageData := &model.RequestToProcess{
		ReqID:    reqID,
		FileName: fileName,
	}

	err = s.processor.ProcessImage(ctx, convertImageData)
	if err != nil {
		return 0, fmt.Errorf("prcoess image: %w", err)
	}

	return reqID, nil
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
	url1, url2, err := s.repo.DeleteRequestAndImage(ctx, userID, reqID)
	if err != nil {
		return err
	}

	err = s.storage.DeleteFile(ctx, url1)
	if err != nil {
		return err
	}

	if url2 != "" {
		err = s.storage.DeleteFile(ctx, url2)
		if err != nil {
			return err
		}
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

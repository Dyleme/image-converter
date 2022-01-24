package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"
)

// ConvertRepo is an interface which provides methods to implement with the reposistory.
type ConvertRepo interface {
	GetConvInfo(ctx context.Context, reqID int) (*model.ConvImageInfo, error)
	SetImageResolution(ctx context.Context, imID int, width int, height int) error
	AddProcessedImage(ctx context.Context, userID, reqID int, imgInfo *model.ReuquestImageInfo,
		width, height int, status string, t time.Time) error
}

// ConvertRequest is a struct provides the ability to convert image.
type ConvertRequest struct {
	repo    ConvertRepo
	storage Storager
	resizer Resizer
}

// NewConvertRequest is a constructor to ConvertRequest.
func NewConvertRequest(repo ConvertRepo, stor Storager, resizer Resizer) *ConvertRequest {
	return &ConvertRequest{repo: repo, storage: stor, resizer: resizer}
}

// Convert is a method to get request from database, image from S3
// convert image and put it all back.
func (c *ConvertRequest) Convert(ctx context.Context, reqID int, filename string) error {
	info, err := c.repo.GetConvInfo(ctx, reqID)
	if err != nil {
		return fmt.Errorf("conversion: %w", err)
	}

	img, err := c.getImage(ctx, info.OldURL, info.OldType)
	if err != nil {
		return fmt.Errorf("conversion: %w", err)
	}

	width, height := getResolution(img)
	err = c.repo.SetImageResolution(ctx, info.OldImID, width, height)

	if err != nil {
		return fmt.Errorf("conversion: %w", err)
	}

	if info.Ratio != 1 {
		img = c.resizer.Resize(img, info.Ratio)
	}

	bts, err := encodeImage(img, info.NewType)
	if err != nil {
		return fmt.Errorf("conversion: %w", err)
	}

	newURL, err := c.storage.UploadFile(ctx, info.UserID, filename, bts)
	if err != nil {
		return fmt.Errorf("conversion: %w", err)
	}

	newImgInfo := model.ReuquestImageInfo{
		URL:  newURL,
		Type: info.NewType,
	}

	newWidth, newHeight := getResolution(img)

	err = c.repo.AddProcessedImage(ctx, info.UserID, reqID, &newImgInfo,
		newWidth, newHeight, repository.StatusDone, time.Now())
	if err != nil {
		return fmt.Errorf("update repo with image: %w", err)
	}

	return nil
}

func (c *ConvertRequest) getImage(ctx context.Context, url, fileType string) (image.Image, error) {
	bts, err := c.storage.GetFile(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("get image: %w", err)
	}

	return decodeImage(bytes.NewBuffer(bts), fileType)
}

package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/Dyleme/image-coverter/internal/conversion"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"
)

var ErrNoPointInFilename = errors.New("no potin in filename")

type NoPointInFilenameError struct {
	Filename string
}

func (e *NoPointInFilenameError) Error() string {
	return fmt.Sprintf("no point in filename: %s", e.Filename)
}

type ConvertRepo interface {
	GetConvInfo(ctx context.Context, reqID int) (*model.ConvImageInfo, error)
	UpdateRequestStatus(ctx context.Context, reqID int, status string) error
	SetImageResolution(ctx context.Context, imID int, width int, height int) error
	AddProcessedImageIDToRequest(ctx context.Context, reqID, imageID int) error
	AddProcessedTimeToRequest(ctx context.Context, reqID int, t time.Time) error
	AddImage(ctx context.Context, userID int, imageInfo model.ReuquestImageInfo) (int, error)
}

type ConvertRequest struct {
	repo    ConvertRepo
	storage Storager
}

func NewConvertRequest(repo ConvertRepo, stor Storager) *ConvertRequest {
	return &ConvertRequest{repo: repo, storage: stor}
}

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
		img = conversion.Resize(img, info.Ratio)
	}

	pointIndex := strings.LastIndex(filename, ".")
	if pointIndex == -1 {
		return &NoPointInFilenameError{Filename: filename}
	}

	convFileName := filename[:pointIndex] + info.NewType

	newImgID, err := c.uploadImage(ctx, img, info.UserID, info.NewType, convFileName)
	if err != nil {
		return fmt.Errorf("conversion: %w", err)
	}

	err = c.updateRepoWithImage(ctx, newImgID, reqID)
	if err != nil {
		return fmt.Errorf("conversion: %w", err)
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

func (c *ConvertRequest) updateRepoWithImage(ctx context.Context, newImgID, reqID int) error {
	err := c.repo.AddProcessedImageIDToRequest(ctx, reqID, newImgID)
	if err != nil {
		return fmt.Errorf("update repo with image: %w", err)
	}

	err = c.repo.AddProcessedTimeToRequest(ctx, reqID, time.Now())
	if err != nil {
		return fmt.Errorf("upodate repo with image: %w", err)
	}

	err = c.repo.UpdateRequestStatus(ctx, reqID, repository.StatusDone)
	if err != nil {
		return fmt.Errorf("update repo with image: %w", err)
	}

	return nil
}

func (c *ConvertRequest) uploadImage(ctx context.Context, img image.Image,
	userID int, fileType, filename string) (int, error) {
	bts, err := encodeImage(img, fileType)
	if err != nil {
		return 0, fmt.Errorf("upload image: %w", err)
	}

	newURL, err := c.storage.UploadFile(ctx, userID, filename, bts)
	if err != nil {
		return 0, fmt.Errorf("upload image: %w", err)
	}

	newImgInfo := model.ReuquestImageInfo{
		URL:  newURL,
		Type: fileType,
	}

	newImgID, err := c.repo.AddImage(ctx, userID, newImgInfo)
	if err != nil {
		return 0, fmt.Errorf("upload image: %w", err)
	}

	newWidth, newHeight := getResolution(img)

	err = c.repo.SetImageResolution(ctx, newImgID, newWidth, newHeight)
	if err != nil {
		return 0, fmt.Errorf("upload image: %w", err)
	}

	return newImgID, nil
}

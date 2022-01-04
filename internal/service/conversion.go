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

func (e NoPointInFilenameError) Error() string {
	return fmt.Sprintf("no point in filename: %s", e.Filename)
}

type ConvertRepo interface {
	GetConvInfo(ctx context.Context, reqID int) (*model.ConvImageInfo, error)
	SetImageResolution(ctx context.Context, imID int, width int, height int) error
	AddImageDB(ctx context.Context, userID, reqID int, imgInfo *model.ReuquestImageInfo,
		width, height int, status string, t time.Time) error
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

	bts, err := encodeImage(img, info.NewType)
	if err != nil {
		return fmt.Errorf("conversion: %w", err)
	}

	newURL, err := c.storage.UploadFile(ctx, info.UserID, convFileName, bts)
	if err != nil {
		return fmt.Errorf("conversion: %w", err)
	}

	newImgInfo := model.ReuquestImageInfo{
		URL:  newURL,
		Type: info.NewType,
	}

	newWidth, newHeight := getResolution(img)

	err = c.repo.AddImageDB(ctx, info.UserID, reqID, &newImgInfo,
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

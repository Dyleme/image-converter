package service

import (
	"context"
	"fmt"

	"github.com/Dyleme/image-coverter/internal/logging"
)

// Download is an interface that provide method gets the image url from the repositoury.
type Downloader interface {
	// GetImageUrl returns the image url.
	GetImageURL(ctx context.Context, userID int, imageID int) (string, error)
}

// DownloadService struct provides the ability to download images from the storage using its id.
type DownloadService struct {
	repo Downloader
	stor Storager
}

// NewDonwloadService is the constructor to the DownloadService.
func NewDownloadService(repo Downloader, stor Storager) *DownloadService {
	return &DownloadService{repo: repo, stor: stor}
}

// Download function returns the bytes of the image or (nil, err) if any error occurs.
// Function gets the imageUrl using repo.GetImgaeURL and get it bytes using stor.GetFile.
func (s *DownloadService) DownloadImage(ctx context.Context, userID, imageID int) ([]byte, error) {
	logger := logging.FromContext(ctx)
	imageURL, err := s.repo.GetImageURL(ctx, userID, imageID)

	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("download image: %w", err)
	}

	fileBytes, err := s.stor.GetFile(ctx, imageURL)

	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("download image: %w", err)
	}

	return fileBytes, nil
}

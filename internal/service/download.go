package service

import (
	"context"
	"fmt"

	"github.com/Dyleme/image-coverter/internal/logging"
)

type Downloader interface {
	GetImageURL(ctx context.Context, userID int, imageID int) (string, error)
}

type DownloadService struct {
	repo Downloader
	stor Storager
}

func NewDownloadSerivce(repo Downloader, stor Storager) *DownloadService {
	return &DownloadService{repo: repo, stor: stor}
}

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

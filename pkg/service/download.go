package service

import (
	"context"
	"fmt"

	"github.com/Dyleme/image-coverter/pkg/repository"
)

type DownloadService struct {
	repo repository.Download
	stor Storager
}

func NewDownloadSerivce(repo repository.Download, stor Storager) *DownloadService {
	return &DownloadService{repo: repo, stor: stor}
}

func (s *DownloadService) DownloadImage(ctx context.Context, userID, imageID int) ([]byte, error) {
	imageURL, err := s.repo.GetImageURL(ctx, userID, imageID)

	if err != nil {
		return nil, fmt.Errorf("download image: %w", err)
	}

	fileBytes, err := s.stor.GetFile(ctx, imageURL)

	if err != nil {
		return nil, fmt.Errorf("download image: %w", err)
	}

	return fileBytes, nil
}

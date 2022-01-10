package service

import (
	"context"
	"fmt"
)

// Download is an interface that provide method gets the image url from the repositoury.
type DownloadRepo interface {
	// GetImageUrl returns the image url.
	GetImageURL(ctx context.Context, userID int, imageID int) (string, error)
}

// Download struct provides the ability to download images from the storage using its id.
type Download struct {
	repo DownloadRepo
	stor Storager
}

// NewDonwloadService is the constructor to the DownloadService.
func NewDownload(repo DownloadRepo, stor Storager) *Download {
	return &Download{repo: repo, stor: stor}
}

// Download function returns the bytes of the image or (nil, err) if any error occurs.
// Function gets the imageUrl using repo.GetImgaeURL and get it bytes using stor.GetFile.
func (s *Download) DownloadImage(ctx context.Context,
	userID, imageID int) (fileBytes []byte, imageURL string, err error) {
	imageURL, err = s.repo.GetImageURL(ctx, userID, imageID)

	if err != nil {
		return nil, "", fmt.Errorf("download image: %w", err)
	}

	fileBytes, err = s.stor.GetFile(ctx, imageURL)

	if err != nil {
		return nil, "", fmt.Errorf("download image: %w", err)
	}

	return fileBytes, imageURL, nil
}

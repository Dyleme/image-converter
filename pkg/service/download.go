package service

import (
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/Dyleme/image-coverter/pkg/storage"
)

type DownloadService struct {
	repo repository.Download
	stor storage.Interface
}

func NewDownloadSerivce(repo repository.Download, stor storage.Interface) *DownloadService {
	return &DownloadService{repo: repo, stor: stor}
}

func (s *DownloadService) DownloadImage(userID, imageID int) ([]byte, error) {
	imageURL, err := s.repo.GetImageURL(userID, imageID)

	if err != nil {
		return nil, err
	}

	fileBytes, err := s.stor.GetFile(imageURL)

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

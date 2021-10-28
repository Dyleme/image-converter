package service

import (
	"github.com/Dyleme/image-coverter"
	"github.com/Dyleme/image-coverter/pkg/repository"
)

type RequestService struct {
	repo repository.Request
}

func NewRequestService(repo repository.Request) *RequestService {
	return &RequestService{repo: repo}
}

func (s *RequestService) GetRequests(userID int) ([]image.Request, error) {
	reqs, err := s.repo.GetRequests(userID)

	if err != nil {
		return nil, err
	}

	return reqs, nil
}

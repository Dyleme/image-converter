package service

import (
	"context"
	"errors"
	"time"

	"github.com/Dyleme/image-coverter/pkg/jwt"
	"github.com/Dyleme/image-coverter/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL = 4 * time.Hour
)

type Autharizater interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetPasswordAndID(ctx context.Context, nickname string) (hash []byte, userID int, err error)
}

type AuthService struct {
	repo Autharizater
}

func NewAuthSevice(repo Autharizater) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(ctx, user)
}

var ErrWrongPassword = errors.New("wrong password")

func (s *AuthService) ValidateUser(ctx context.Context, user model.User) (string, error) {
	hash, id, err := s.repo.GetPasswordAndID(ctx, user.Nickname)
	if err != nil {
		return "", ErrWrongPassword
	}

	if !isValidPassword(user.Password, hash) {
		return "", ErrWrongPassword
	}

	return jwt.CreateToketn(ctx, tokenTTL, id)
}

func generatePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

func isValidPassword(password string, hash []byte) bool {
	errNotEqual := bcrypt.CompareHashAndPassword(hash, []byte(password))

	return errNotEqual == nil
}

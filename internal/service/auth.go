package service

import (
	"context"
	"errors"
	"time"

	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/Dyleme/image-coverter/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL = 4 * time.Hour
)

type HashGenerator interface {
	GeneratePasswordHash(password string) string
	IsValidPassword(password string, hash []byte) bool
}

type HashGen struct{}

func (h *HashGen) GeneratePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

func (h *HashGen) IsValidPassword(password string, hash []byte) bool {
	errNotEqual := bcrypt.CompareHashAndPassword(hash, []byte(password))

	return errNotEqual == nil
}

type Autharizater interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetPasswordAndID(ctx context.Context, nickname string) (hash []byte, userID int, err error)
}

type JwtGen struct{}

func NewJwtGen() *JwtGen {
	return &JwtGen{}
}

func (gen *JwtGen) CreateToken(ctx context.Context, tokenTTL time.Duration, id int) (string, error) {
	return jwt.CreateToken(ctx, tokenTTL, id)
}

type AuthService struct {
	repo    Autharizater
	hashGen HashGenerator
	jwtGen  JwtGenerator
}

func NewAuthSevice(repo Autharizater, hashGen HashGenerator, jwtGen JwtGenerator) *AuthService {
	return &AuthService{repo: repo, hashGen: hashGen, jwtGen: jwtGen}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (int, error) {
	user.Password = s.hashGen.GeneratePasswordHash(user.Password)
	return s.repo.CreateUser(ctx, user)
}

var ErrWrongPassword = errors.New("wrong password")

func (s *AuthService) ValidateUser(ctx context.Context, user model.User) (string, error) {
	hash, id, err := s.repo.GetPasswordAndID(ctx, user.Nickname)
	if err != nil {
		return "", ErrWrongPassword
	}

	if !s.hashGen.IsValidPassword(user.Password, hash) {
		return "", ErrWrongPassword
	}

	return s.jwtGen.CreateToken(ctx, tokenTTL, id)
}

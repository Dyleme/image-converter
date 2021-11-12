package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Dyleme/image-coverter/pkg/model"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL  = 4 * time.Hour
	signedKey = "2lkj^@dkjg#)jfkdlg"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthSevice(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(ctx, user)
}

var ErrWrongPassword = errors.New("wrong password")

var ErrTokenClaimsInvalidType = errors.New("token claims are not of the type MapClaims")

type UnexpectedSingingMethodError struct {
	method interface{}
}

func (err UnexpectedSingingMethodError) Error() string {
	return fmt.Sprintf("unexpected singing method: %v", err.method)
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"UserID"`
}

func (s *AuthService) ValidateUser(ctx context.Context, user model.User) (string, error) {
	hash, id, err := s.repo.GetPasswordAndID(ctx, user.Nickname)
	if err != nil {
		return "", ErrWrongPassword
	}

	if !isValidPassword(user.Password, hash) {
		return "", ErrWrongPassword
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		},
		id,
	})

	return jwtToken.SignedString([]byte(signedKey))
}

func (s *AuthService) ParseToken(ctx context.Context, tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, UnexpectedSingingMethodError{t.Header["alg"]}
		}

		return []byte(signedKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, err := strconv.Atoi(fmt.Sprintf("%.f", claims["UserID"]))
		if err != nil {
			return 0, fmt.Errorf("parse token: %w", err)
		}

		return userID, nil
	}

	return 0, ErrTokenClaimsInvalidType
}

func generatePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

func isValidPassword(password string, hash []byte) bool {
	errNotEqual := bcrypt.CompareHashAndPassword(hash, []byte(password))

	return errNotEqual == nil
}

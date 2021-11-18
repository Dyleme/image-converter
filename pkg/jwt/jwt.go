package jwt

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type Key string

var signedKey = os.Getenv("SIGNEDKEY")

const (
	KeyUserID Key = "keyUserID"
)

var ErrTokenClaimsInvalidType = errors.New("token claims are not of the type MapClaims")

var ErrContextHaveNotUser = errors.New("can't get user from context")

type UnexpectedSingingMethodError struct {
	method interface{}
}

func (err UnexpectedSingingMethodError) Error() string {
	return fmt.Sprintf("unexpected singing method: %v", err.method)
}

type tokenClaims struct {
	jwt.Claims
	UserID int `json:"userID"`
}

func CreateToken(ctx context.Context, tokenTTL time.Duration, id int) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		},
		id,
	})

	return jwtToken.SignedString([]byte(signedKey))
}

func ParseToken(ctx context.Context, tokenString string) (int, error) {
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
		userID, err := strconv.Atoi(fmt.Sprintf("%.f", claims["userID"]))
		if err != nil {
			return 0, fmt.Errorf("parse token: %w", err)
		}

		return userID, nil
	}

	return 0, ErrTokenClaimsInvalidType
}

func GetUserFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(KeyUserID).(int)

	if !ok {
		return 0, ErrContextHaveNotUser
	}

	return userID, nil
}

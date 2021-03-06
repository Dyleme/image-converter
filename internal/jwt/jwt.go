package jwt

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type Key string

type Gen struct {
	signedKey string
	ttl       time.Duration
}

type Config struct {
	SignedKey string
	TTL       time.Duration
}

func NewJwtGen(config *Config) *Gen {
	return &Gen{signedKey: config.SignedKey}
}

const (
	KeyUserID Key = "keyUserID"
)

var ErrTokenClaimsInvalidType = errors.New("token claims are not of the type MapClaims")

var ErrContextWithoutUser = errors.New("can't get user from context")

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

// CreateToken function generate token with provided TTL and user id.
func (g *Gen) CreateToken(_ context.Context, id int) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(g.ttl).Unix(),
		},
		id,
	})

	return jwtToken.SignedString([]byte(g.signedKey))
}

// ParseToken function rerurns user id from JWT token, if this token is liquid.
func (g *Gen) ParseToken(_ context.Context, tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, UnexpectedSingingMethodError{t.Header["alg"]}
		}

		return []byte(g.signedKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("parse token: %w", err)
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

// Function GetUserFromContext return a user from a context,
// or error ErrContextWithoutUser if it isn't user in context.
func GetUserFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(KeyUserID).(int)

	if !ok {
		return 0, ErrContextWithoutUser
	}

	return userID, nil
}

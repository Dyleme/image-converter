package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type key string

const (
	AuthorizationHeader = "Authorization"

	BearerToken = "Bearer"

	keyUserID key = "keyUserID"
)

var ErrContextHaveNotUser = errors.New("can't get user from context")

func (h *Handler) checkJWT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader, exist := r.Header[AuthorizationHeader]
		if !exist {
			newErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		if len(authHeader) != 1 {
			newErrorResponse(w, http.StatusUnauthorized, "more than one auth header")
			return
		}

		auth := authHeader[0]

		if auth[:len(BearerToken)] != BearerToken {
			newErrorResponse(w, http.StatusUnauthorized, "invalid authentication method")
			return
		}

		authJWT := auth[len(BearerToken):]
		authJWT = strings.TrimPrefix(authJWT, " ")

		userID, err := h.service.ParseToken(ctx, authJWT)
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("middleware: %w", err).Error())
			return
		}

		ctx = context.WithValue(ctx, keyUserID, userID)

		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("request %v method %v time for answer : %v", r.URL.Path, r.Method, time.Since(begin))
	})
}

func getUserFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(keyUserID).(int)

	if !ok {
		return 0, ErrContextHaveNotUser
	}

	return userID, nil
}

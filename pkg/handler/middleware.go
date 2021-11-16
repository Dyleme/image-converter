package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Dyleme/image-coverter/pkg/jwt"
)

const (
	AuthorizationHeader = "Authorization"

	BearerToken = "Bearer"
)

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

		userID, err := jwt.ParseToken(ctx, authJWT)
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("middleware: %w", err).Error())
			return
		}

		ctx = context.WithValue(ctx, jwt.KeyUserID, userID)

		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		handler.ServeHTTP(w, r)
		file, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeAppend)
		if err != nil {
			log.Fatal()
		}
		log.SetOutput(file)
		log.WithFields(log.Fields{
			"request path":    r.URL.Path,
			"reqest method":   r.Method,
			"time for answer": time.Since(begin),
		}).Info("request handled")
	})
}

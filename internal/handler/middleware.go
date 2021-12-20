package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Dyleme/image-coverter/internal/jwt"
)

const (
	AuthorizationHeader = "Authorization"

	BearerToken = "Bearer"
)

// checkJWT is not allow to acces to endpoints to unauthorized users.
// It get jwt tokent from a header, get the user by this token and put the user to the context.
// If eroor doesn't occur than the next hadnler handle the reqeust,
// else method response with the error answer.
func CheckJWT(handler http.Handler) http.Handler {
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

// logging is used to log all incoming requests.
// It logs request's url, method and time for answer to this reqeust.
func (h *Handler) logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		h.logger.WithFields(log.Fields{
			"request path":    r.URL.Path,
			"reqest method":   r.Method,
			"time for answer": begin,
		}).Info("get request")

		handler.ServeHTTP(w, r)

		h.logger.WithFields(log.Fields{
			"request path":    r.URL.Path,
			"reqest method":   r.Method,
			"time for answer": time.Since(begin),
		}).Info("request handled")
	})
}

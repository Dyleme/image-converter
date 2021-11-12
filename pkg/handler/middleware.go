package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type key string

const (
	AuthorizationHeader = "Authorization"

	keyUserID key = "keyUserID"
)

const (
	BearerToken = "Bearer"
)

func (h *Handler) checkJWT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader, exist := r.Header[AuthorizationHeader]
		if !exist {
			newErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(authHeader[0], " ")
		log.Println(headerParts[0])

		if len(headerParts) != 2 { //nolint:gomnd // 2 is amount of argumetns that should have auth
			newErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		userID, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := r.Context()

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

func getUserFromContext(r *http.Request) (int, error) {
	ctx := r.Context()
	userID, ok := ctx.Value(keyUserID).(int)

	if !ok {
		return 0, fmt.Errorf("can't get user from context")
	}

	return userID, nil
}

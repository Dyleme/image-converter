package controller

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

type key int

const (
	AuthorizationHeader = "Authorization"

	keyUserID key = iota
)

func (c *Controller) checkJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader, exist := r.Header[AuthorizationHeader]
		if !exist {
			newErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(authHeader[0], " ")

		if len(headerParts) != 2 { //nolint:gomnd // 2 is amount of argumetns that should have auth
			newErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		userID, err := c.service.ParseToken(headerParts[1])
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(context.Background(), keyUserID, userID)

		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

func (c *Controller) log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("time for answer : %v", time.Since(begin))
	})
}

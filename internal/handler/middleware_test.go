package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Dyleme/image-coverter/internal/handler"
	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/stretchr/testify/assert"
)

type handlerMock struct{}

func (hm *handlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := jwt.GetUserFromContext(r.Context())
	if err != nil {
		fmt.Fprint(w, []byte(err.Error()))
	}

	fmt.Fprint(w, strconv.Itoa(id))
}

var handMock handlerMock

func TestCheckJwt(t *testing.T) {
	testCases := []struct {
		testName     string
		method       string
		path         string
		reqHeaderKey string
		reqHeaderVal string
		initHeader   func(*http.Request) *handler.JwtChecker
		wantStatus   int
		wantBody     string
	}{
		{
			testName:     "ok",
			reqHeaderKey: "Authorization",
			initHeader: func(req *http.Request) *handler.JwtChecker {
				checker := &handler.JwtChecker{*jwt.NewJwtGen(&jwt.Config{
					SignedKey: "key",
					TTL:       4 * time.Hour,
				})}
				token, _ := checker.Gen.CreateToken(context.Background(), 12)
				req.Header.Add("Authorization", "Bearer "+token)
				return checker
			},
			reqHeaderVal: "Bearer",
			wantStatus:   http.StatusOK,
			wantBody:     "12",
		},
		{
			testName: "empty auth header",
			initHeader: func(req *http.Request) *handler.JwtChecker {
				checker := &handler.JwtChecker{*jwt.NewJwtGen(&jwt.Config{
					SignedKey: "key",
					TTL:       4 * time.Hour,
				})}
				return checker
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"message":"empty auth header"}`,
		},
		{
			testName: "multiply auth headers",
			initHeader: func(req *http.Request) *handler.JwtChecker {
				checker := &handler.JwtChecker{*jwt.NewJwtGen(&jwt.Config{
					SignedKey: "key",
					TTL:       4 * time.Hour,
				})}
				token, _ := checker.Gen.CreateToken(context.Background(), 12)
				req.Header.Add("Authorization", "Bearer "+token)
				req.Header.Add("Authorization", "Bearer "+token)
				return checker
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"message":"more than one auth header"}`,
		},
		{
			testName: "invalide auth method",
			initHeader: func(req *http.Request) *handler.JwtChecker {
				checker := handler.JwtChecker{*jwt.NewJwtGen(&jwt.Config{
					SignedKey: "key",
					TTL:       4 * time.Hour,
				})}
				token, _ := checker.Gen.CreateToken(context.Background(), 12)
				req.Header.Add("Authorization", "Invalide "+token)
				return &checker
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"message":"invalid authentication method"}`,
		},
		{
			testName: "invalide jwt token",
			initHeader: func(req *http.Request) *handler.JwtChecker {
				checker := handler.JwtChecker{*jwt.NewJwtGen(&jwt.Config{
					SignedKey: "key",
					TTL:       4 * time.Hour,
				})}
				token, _ := checker.Gen.CreateToken(context.Background(), 12)
				req.Header.Add("Authorization", "Bearer "+token+"to invalid")
				return &checker
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"message":"middleware: parse token: illegal base64 data at input byte 45"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, &strings.Reader{})
			if err != nil {
				t.Fatal(err)
			}
			checker := tc.initHeader(req)

			rr := httptest.NewRecorder()

			checker.CheckJWT(&handMock).ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tc.wantStatus)
			assert.Equal(t, rr.Body.String(), tc.wantBody)
		})
	}
}

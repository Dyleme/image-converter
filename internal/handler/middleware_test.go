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
		initHeader   func(*http.Request)
		wantStatus   int
		wantBody     string
	}{
		{
			testName:     "ok",
			reqHeaderKey: "Authorization",
			initHeader: func(req *http.Request) {
				token, _ := jwt.CreateToken(context.Background(), time.Hour, 12)
				req.Header.Add("Authorization", "Bearer "+token)
			},
			reqHeaderVal: "Bearer",
			wantStatus:   http.StatusOK,
			wantBody:     "12",
		},
		{
			testName: "empty auth header",
			initHeader: func(req *http.Request) {
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"message":"empty auth header"}`,
		},
		{
			testName: "multiply auth headers",
			initHeader: func(req *http.Request) {
				token, _ := jwt.CreateToken(context.Background(), time.Hour, 12)
				req.Header.Add("Authorization", "Bearer "+token)
				req.Header.Add("Authorization", "Bearer "+token)
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"message":"more than one auth header"}`,
		},
		{
			testName: "invalide auth method",
			initHeader: func(req *http.Request) {
				token, _ := jwt.CreateToken(context.Background(), time.Hour, 12)
				req.Header.Add("Authorization", "Invalide "+token)
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"message":"invalid authentication method"}`,
		},
		{
			testName: "invalide jwt token",
			initHeader: func(req *http.Request) {
				token, _ := jwt.CreateToken(context.Background(), time.Hour, 12)
				req.Header.Add("Authorization", "Bearer "+token+"to invalid")
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

			tc.initHeader(req)

			rr := httptest.NewRecorder()

			handler.CheckJWT(&handMock).ServeHTTP(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("want status %v, got status %v", tc.wantStatus, status)
			}

			if body := rr.Body.String(); body != tc.wantBody {
				t.Errorf("want body %s, got body %s", tc.wantBody, body)
			}
		})
	}
}

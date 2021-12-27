package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dyleme/image-coverter/internal/handler"
	"github.com/Dyleme/image-coverter/internal/handler/mocks"
	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var errDownloading = errors.New("error in downloading")

func Test(t *testing.T) {
	testCases := []struct {
		testName   string
		method     string
		path       string
		imageID    string
		configure  func(*http.Request, *mocks.MockDownloader) *http.Request
		wantStatus int
		wantBody   string
	}{
		{
			testName: "ok",
			method:   http.MethodGet,
			path:     "image/download/12",
			configure: func(r *http.Request, md *mocks.MockDownloader) *http.Request {
				md.EXPECT().DownloadImage(gomock.Any(), 2, 12).Return([]byte("body"), nil).Times(1)

				r = mux.SetURLVars(r, map[string]string{
					"id": "12",
				})

				ctx := context.WithValue(r.Context(), jwt.KeyUserID, 2)

				return r.WithContext(ctx)
			},
			wantStatus: http.StatusOK,
			wantBody:   "body",
		},
		{
			testName: "no auth",
			method:   http.MethodGet,
			path:     "image/download/12",
			configure: func(r *http.Request, md *mocks.MockDownloader) *http.Request {
				return r
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"message":"can't get user from context"}`,
		},
		{
			testName: "parameter is missing",
			method:   http.MethodGet,
			path:     "image/download/",
			configure: func(r *http.Request, md *mocks.MockDownloader) *http.Request {
				ctx := context.WithValue(r.Context(), jwt.KeyUserID, 2)
				return r.WithContext(ctx)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"parameter \"id\" is missing"}`,
		},
		{
			testName: "other parameter is provided",
			method:   http.MethodGet,
			path:     "image/download/a",
			configure: func(r *http.Request, md *mocks.MockDownloader) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"a": "12",
				})
				ctx := context.WithValue(r.Context(), jwt.KeyUserID, 2)
				return r.WithContext(ctx)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"parameter \"id\" is missing"}`,
		},
		{
			testName: "id is not int",
			method:   http.MethodGet,
			path:     "image/download/a",
			configure: func(r *http.Request, md *mocks.MockDownloader) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"id": "not int",
				})
				ctx := context.WithValue(r.Context(), jwt.KeyUserID, 2)
				return r.WithContext(ctx)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"strconv.Atoi: parsing \"not int\": invalid syntax"}`,
		},
		{
			testName: "ok",
			method:   http.MethodGet,
			path:     "image/download/12",
			configure: func(r *http.Request, md *mocks.MockDownloader) *http.Request {
				md.EXPECT().DownloadImage(gomock.Any(), 2, 12).Return(nil, errDownloading).Times(1)

				r = mux.SetURLVars(r, map[string]string{
					"id": "12",
				})

				ctx := context.WithValue(r.Context(), jwt.KeyUserID, 2)

				return r.WithContext(ctx)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"error in downloading"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			req, err := http.NewRequest(tc.method, tc.path, &strings.Reader{})
			if err != nil {
				t.Fatal(err)
			}

			downMock := mocks.NewMockDownloader(mockCtr)
			downHandler := handler.NewDownload(downMock, &logrus.Logger{})

			req = tc.configure(req, downMock)

			rr := httptest.NewRecorder()

			downHandler.DownloadImage(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("want status %v, got status %v", tc.wantStatus, status)
			}

			if body := rr.Body.String(); body != tc.wantBody {
				t.Errorf("want body %s, got body %s", tc.wantBody, body)
			}
		})
	}
}

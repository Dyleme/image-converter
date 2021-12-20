package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dyleme/image-coverter/internal/handler"
	"github.com/Dyleme/image-coverter/internal/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
)

var errInvalidUser = errors.New("invalid user")

func TestAuthHandler_Login(t *testing.T) {
	testCases := []struct {
		testName   string
		method     string
		path       string
		reqBody    string
		configure  func(*mocks.MockAutharizater)
		wantStatus int
		wantBody   string
	}{
		{
			testName:   "empty request body",
			method:     http.MethodGet,
			path:       "/auth/login",
			configure:  func(ma *mocks.MockAutharizater) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"EOF"}`,
		},
		{
			testName:  "unmarshaling request body",
			method:    http.MethodGet,
			path:      "/auth/login",
			configure: func(ma *mocks.MockAutharizater) {},
			reqBody: `{
    			"Nickname : "Dyleme",
  		  		"Password": "1231"
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"invalid character 'D' after object key"}`,
		},
		{
			testName:  "body not full",
			method:    http.MethodGet,
			path:      "/auth/login",
			configure: func(ma *mocks.MockAutharizater) {},
			reqBody: `{
    			"Nickname" : "Dyleme"
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"all fields should be filled {Dyleme,}"}`,
		},
		{
			testName: "all is good",
			method:   http.MethodGet,
			path:     "/auth/login",
			configure: func(ma *mocks.MockAutharizater) {
				ma.EXPECT().ValidateUser(gomock.Any(), gomock.Any()).Return("jwt", nil).Times(1)
			},
			reqBody: `{
    			"Nickname": "Dyleme",
  		  		"Password": "1231"
			}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"jwt":"jwt"}`,
		},
		{
			testName: "invalid user",
			method:   http.MethodGet,
			path:     "/auth/login",
			configure: func(ma *mocks.MockAutharizater) {
				ma.EXPECT().ValidateUser(gomock.Any(), gomock.Any()).Return("", errInvalidUser).Times(1)
			},
			reqBody: `{
    			"Nickname": "Dyleme",
  		  		"Password": "1231"
			}`,
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"invalid user"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			req, err := http.NewRequest(tc.method, tc.path, strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			authMock := mocks.NewMockAutharizater(mockCtr)
			authHandler := handler.NewAuthHandler(authMock, &logrus.Logger{})

			tc.configure(authMock)

			authHandler.Login(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("want status %v, got status %v", tc.wantStatus, status)
			}

			if body := rr.Body.String(); body != tc.wantBody {
				t.Errorf("want body %s, got body %s", tc.wantBody, body)
			}
		})
	}
}

var errCreatingUser = errors.New("creating user error")

func TestAuthHandler_Register(t *testing.T) {
	testCases := []struct {
		testName   string
		method     string
		path       string
		reqBody    string
		configure  func(*mocks.MockAutharizater)
		wantStatus int
		wantBody   string
	}{
		{
			testName: "all is good",
			method:   http.MethodGet,
			path:     "/auth/register",
			configure: func(ma *mocks.MockAutharizater) {
				ma.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(1, nil).Times(1)
			},
			reqBody: `{
    			"Nickname": "Dyleme",
  		  		"Password": "1231"
			}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"id":1}`,
		},
		{
			testName: "error in creating user",
			method:   http.MethodGet,
			path:     "/auth/register",
			configure: func(ma *mocks.MockAutharizater) {
				ma.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(0, errCreatingUser).Times(1)
			},
			reqBody: `{
    			"Nickname": "Dyleme",
  		  		"Password": "1231"
			}`,
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"creating user error"}`,
		},
		{
			testName:   "empty request body",
			method:     http.MethodGet,
			path:       "/auth/register",
			configure:  func(ma *mocks.MockAutharizater) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"EOF"}`,
		},
		{
			testName:  "unmarshaling request body",
			method:    http.MethodGet,
			path:      "/auth/register",
			configure: func(ma *mocks.MockAutharizater) {},
			reqBody: `{
    			"Nickname : "pyleme",
  		  		"Password": "1231"
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"invalid character 'p' after object key"}`,
		},
		{
			testName:  "body not full",
			method:    http.MethodGet,
			path:      "/auth/register",
			configure: func(ma *mocks.MockAutharizater) {},
			reqBody: `{
    			"Nickname" : "Dyleme"
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"all fields should be filled {Dyleme,}"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			req, err := http.NewRequest(tc.method, tc.path, strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			authMock := mocks.NewMockAutharizater(mockCtr)

			tc.configure(authMock)

			authHandler := handler.NewAuthHandler(authMock, &logrus.Logger{})

			authHandler.Register(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("want status %v, got status %v", tc.wantStatus, status)
			}

			if body := rr.Body.String(); body != tc.wantBody {
				t.Errorf("want body %s, got body %s", tc.wantBody, body)
			}
		})
	}
}

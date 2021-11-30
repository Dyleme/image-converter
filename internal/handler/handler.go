package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	authService     Autharizater
	requestService  Requester
	downloadService Downloader
	logger          *logrus.Logger
}

func New(auth Autharizater, request Requester, download Downloader, logger *logrus.Logger) *Handler {
	return &Handler{authService: auth, requestService: request, downloadService: download, logger: logger}
}

type Autharizater interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	ValidateUser(ctx context.Context, user model.User) (string, error)
}

type Requester interface {
	GetRequests(ctx context.Context, userID int) ([]model.Request, error)
	GetRequest(ctx context.Context, userID int, reqID int) (*model.Request, error)
	DeleteRequest(ctx context.Context, userID int, reqID int) error
	AddRequest(context.Context, int, io.Reader, string, model.ConversionInfo) (int, error)
}

type Downloader interface {
	DownloadImage(ctx context.Context, userID, imageID int) ([]byte, error)
}

func (h *Handler) InitRouters() *mux.Router {
	router := mux.NewRouter()

	authRouter := router.NewRoute().Subrouter()

	router.HandleFunc("/auth/register", h.RegiterHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", h.LoginHandler)

	authRouter.HandleFunc("/requests", h.AllRequestsHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/requests/image", h.AddRequestHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/requests/{reqID}", h.GetRequestHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/requests/{reqID}", h.DeleteRequestHandler).Methods(http.MethodDelete)

	authRouter.HandleFunc("/download/image/{id}", h.DownloadImageHandler).Methods(http.MethodGet)

	router.Use(h.logging)
	authRouter.Use(h.checkJWT)

	return router
}
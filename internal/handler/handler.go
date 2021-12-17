package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Handler is a struct which has service interfaces.
type Handler struct {
	authService     Autharizater
	requestService  Requester
	downloadService Downloader

	authHandler AuthenticationHandler
	reqHandler  RequestHandler
	downHandler DownloadHandler

	// logger is used to write all logs in Handler
	logger *logrus.Logger
}

// This constructor initialize Handler's fields with provided arguments.
func New(authHand AuthenticationHandler, reqHandler RequestHandler, downHandler DownloadHandler,
	authServ Autharizater, reqServ Requester, downServ Downloader, logger *logrus.Logger) *Handler {
	return &Handler{authHandler: authHand, reqHandler: reqHandler, downHandler: downHandler,
		authService: authServ, requestService: reqServ, downloadService: downServ, logger: logger}
}

// Autharizater is an interface which has methods to create and validate user.
type Autharizater interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	ValidateUser(ctx context.Context, user model.User) (string, error)
}

// Requester is an interface which has methods to get, delete and add requests.
type Requester interface {
	GetRequests(ctx context.Context, userID int) ([]model.Request, error)
	GetRequest(ctx context.Context, userID int, reqID int) (*model.Request, error)
	DeleteRequest(ctx context.Context, userID int, reqID int) error
	AddRequest(context.Context, int, io.Reader, string, model.ConversionInfo) (int, error)
}

// Downloader is an interface which has method to download image.
type Downloader interface {
	DownloadImage(ctx context.Context, userID, imageID int) ([]byte, error)
}

type AuthenticationHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type H interface {
	AuthenticationHandler
	DownloadHandler
}

type RequestHandler interface {
	GetAllRequests(w http.ResponseWriter, r *http.Request)
	GetRequest(w http.ResponseWriter, r *http.Request)
	AddRequest(w http.ResponseWriter, r *http.Request)
	DeleteRequest(w http.ResponseWriter, r *http.Request)
}

type DownloadHandler interface {
	DownloadImage(w http.ResponseWriter, r *http.Request)
}

// InitRouters() method is used to initialize all endopoints with the routers.
func (h *Handler) InitRouters() *mux.Router {
	router := mux.NewRouter()

	authRouter := router.NewRoute().Subrouter()

	router.HandleFunc("/auth/register", h.authHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", h.authHandler.Login).Methods(http.MethodPost)

	authRouter.HandleFunc("/requests", h.reqHandler.GetAllRequests).Methods(http.MethodGet)
	authRouter.HandleFunc("/requests/{reqID}", h.reqHandler.GetRequest).Methods(http.MethodGet)
	authRouter.HandleFunc("/requests/image", h.reqHandler.AddRequest).Methods(http.MethodPost)
	authRouter.HandleFunc("/requests/{reqID}", h.reqHandler.DeleteRequest).Methods(http.MethodDelete)

	authRouter.HandleFunc("/download/image/{id}", h.downHandler.DownloadImage).Methods(http.MethodGet)

	router.Use(h.logging)
	authRouter.Use(h.checkJWT)

	return router
}

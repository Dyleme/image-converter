package handler

import (
	"net/http"

	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Handler is a struct which has service interfaces.
type Handler struct {
	authHandler AuthenticationHandler
	reqHandler  RequestHandler
	downHandler DownloadHandler

	// logger is used to write all logs in Handler
	logger *logrus.Logger
}

// This constructor initialize Handler's fields with provided arguments.
func New(authHand AuthenticationHandler, reqHandler RequestHandler, downHandler DownloadHandler,
	logger *logrus.Logger) *Handler {
	return &Handler{authHandler: authHand, reqHandler: reqHandler, downHandler: downHandler,
		logger: logger}
}

type ErrorWithStatus interface {
	error
	Status() int
}

type AuthenticationHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
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
func (h *Handler) InitRouters(jwtGen *jwt.Gen) *mux.Router {
	router := mux.NewRouter()
	router.Use(h.logTime)

	authRouter := router.NewRoute().Subrouter()
	jwtChecker := JwtChecker{Gen: *jwtGen}
	authRouter.Use(jwtChecker.CheckJWT)

	router.HandleFunc("/auth/register", h.authHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", h.authHandler.Login).Methods(http.MethodPost)

	authRouter.HandleFunc("/requests", h.reqHandler.GetAllRequests).Methods(http.MethodGet)
	authRouter.HandleFunc("/requests/{reqID}", h.reqHandler.GetRequest).Methods(http.MethodGet)
	authRouter.HandleFunc("/requests/image", h.reqHandler.AddRequest).Methods(http.MethodPost)
	authRouter.HandleFunc("/requests/{reqID}", h.reqHandler.DeleteRequest).Methods(http.MethodDelete)

	authRouter.HandleFunc("/download/image/{id}", h.downHandler.DownloadImage).Methods(http.MethodGet)

	return router
}

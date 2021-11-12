package handler

import (
	"net/http"

	"github.com/Dyleme/image-coverter/pkg/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service service.Interface
}

func NewServer(serv service.Interface) *Handler {
	return &Handler{service: serv}
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

	router.Use(logging)
	authRouter.Use(h.checkJWT)

	return router
}

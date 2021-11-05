package handler

import (
	"net/http"

	"github.com/Dyleme/image-coverter/pkg/service"
	"github.com/gorilla/mux"
)

type Server struct {
	service service.Interface
}

func NewController(serv service.Interface) *Server {
	return &Server{service: serv}
}

func (s *Server) InitRouters() *mux.Router {
	router := mux.NewRouter()

	authRouter := router.NewRoute().Subrouter()

	router.HandleFunc("/auth/register", s.RegiterHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", s.LoginHandler)

	authRouter.HandleFunc("/requests", s.AllRequestsHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/requests/image", s.AddRequestHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/requests/{reqID}", s.GetRequestHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/requests/{reqID}", s.DeleteRequestHandler).Methods(http.MethodDelete)

	authRouter.HandleFunc("/download/image/{id}", s.DownloadImageHandler).Methods(http.MethodGet)

	router.Use(s.log)
	authRouter.Use(s.checkJWT)

	return router
}

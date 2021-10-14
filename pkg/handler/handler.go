package handler

import (
	"github.com/gorilla/mux"

	"github.com/Dyleme/image-coverter/pkg/service"
)

type Handler struct {
	service *service.Service
}

func (h *Handler) InitRouters() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/auth/sign-in", SignInHandler)
	r.HandleFunc("/auth/sign-up", SignUpHandler)

	r.HandleFunc("/history", HistortHandler)

	r.HandleFunc("/conversion", ConversionHandler)
	return r
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

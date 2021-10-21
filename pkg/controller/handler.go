package controller

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Dyleme/image-coverter/pkg/service"
)

type Controller struct {
	service *service.Service
}

func NewController(serv *service.Service) *Controller {
	return &Controller{service: serv}
}

func (c *Controller) InitRouters() *mux.Router {
	router := mux.NewRouter()

	authRouter := router.NewRoute().Subrouter()

	router.HandleFunc("/auth/register", c.RegiterHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", c.LoginHandler)

	authRouter.HandleFunc("/requests", c.AllRequestsHandler).Methods(http.MethodGet)

	authRouter.Use(c.log)
	authRouter.Use(c.checkJWT)

	return router
}

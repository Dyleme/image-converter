package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/sirupsen/logrus"
)

// Struct which provides methods to handle Login and Registration.
type AuthHandler struct {
	logger      *logrus.Logger
	authService Autharizater
}

// Constructor for AuthHandler.
func NewAuthHandler(auth Autharizater, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{authService: auth, logger: logger}
}

// Login is method which decode request body to model.User
// And validate this user using ValidateUser service method.
// Method response with jwt token or error, if any occurs.
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		ah.logger.Warn(err)

		return
	}

	jwtToken, err := ah.authService.ValidateUser(ctx, input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		ah.logger.Warn(err)

		return
	}

	newJSONResponse(w, jwtToken)
}

// Register is function which is used to register users.
// It decodes request body to model.User
// and calls service method CreateUser.
// Method response with user id or error, if any occurs.
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ah.logger.Warn(err)
		newErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	id, err := ah.authService.CreateUser(ctx, input)
	if err != nil {
		ah.logger.Warn(err)
		newErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	idStruct := map[string]int{"id": id}
	newJSONResponse(w, idStruct)
}

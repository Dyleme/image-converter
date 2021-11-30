package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Dyleme/image-coverter/internal/model"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		h.logger.Warn(err)

		return
	}

	jwtToken, err := h.authService.ValidateUser(ctx, input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		h.logger.Warn(err)

		return
	}

	newJSONResponse(w, jwtToken)
}

func (h *Handler) RegiterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Warn(err)
		newErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	id, err := h.authService.CreateUser(ctx, input)
	if err != nil {
		h.logger.Warn(err)
		newErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	idStruct := map[string]int{"id": id}
	newJSONResponse(w, idStruct)
}
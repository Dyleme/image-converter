package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Dyleme/image-coverter/pkg/model"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := h.service.ValidateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, jwtToken)
}

func (h *Handler) RegiterHandler(w http.ResponseWriter, r *http.Request) {
	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	idStruct := map[string]int{"id": id}
	newJSONResponse(w, idStruct)
}

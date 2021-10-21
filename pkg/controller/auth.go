package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Dyleme/image-coverter"
)

func (c *Controller) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input image.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := c.service.ValidateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, jwtToken)
}

func (c *Controller) RegiterHandler(w http.ResponseWriter, r *http.Request) {
	var input image.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.service.CreateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	idStruct := map[string]int{"id": id}
	jsonResponse(w, idStruct)
}

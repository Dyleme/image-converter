package controller

import (
	"net/http"
)

func (c *Controller) AllRequestsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := ctx.Value(keyUserID).(int)

	if !ok {
		newErrorResponse(w, http.StatusUnauthorized, "wrong user id")
		return
	}

	reqs, err := c.service.GetRequests(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, reqs)
}

func (c *Controller) AddRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := ctx.Value(keyUserID).(int)

	if !ok {
		newErrorResponse(w, http.StatusUnauthorized, "wrong user id came from middleware")
		return
	}

	reqs, err := c.service.GetRequests(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, reqs)
}

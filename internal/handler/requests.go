package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/gorilla/mux"
)

func (h *Handler) AllRequestsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := jwt.GetUserFromContext(ctx)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	reqs, err := h.requestService.GetRequests(ctx, userID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, reqs)
}

func (h *Handler) AddRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := jwt.GetUserFromContext(ctx)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	file, header, err := r.FormFile("Image")
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer file.Close()

	info := r.FormValue("CompressionInfo")

	var sendInfo model.ConversionInfo

	err = json.Unmarshal([]byte(info), &sendInfo)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	reqID, err := h.requestService.AddRequest(ctx, userID, file, header.Filename, sendInfo)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	m := struct {
		RequestID int `json:"requestID"`
	}{
		RequestID: reqID,
	}

	newJSONResponse(w, m)
}

func (h *Handler) GetRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := jwt.GetUserFromContext(ctx)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)

	strReqID, ok := vars["reqID"]
	if !ok {
		newErrorResponse(w, http.StatusBadRequest, "id parameter is missing")
		return
	}

	reqID, err := strconv.Atoi(strReqID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	request, err := h.requestService.GetRequest(ctx, userID, reqID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, request)
}

func (h *Handler) DeleteRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := jwt.GetUserFromContext(ctx)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)

	strReqID, ok := vars["reqID"]
	if !ok {
		newErrorResponse(w, http.StatusBadRequest, "id parameter is missing")
		return
	}

	reqID, err := strconv.Atoi(strReqID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.requestService.DeleteRequest(ctx, userID, reqID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, reqID)
}

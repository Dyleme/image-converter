package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Struct which provides methods to handle working with requests.
type ReqHandler struct {
	logger         *logrus.Logger
	requestService Requester
}

// Constructor for ReqHandler.
func NewReqHandler(req Requester, logger *logrus.Logger) *ReqHandler {
	return &ReqHandler{requestService: req, logger: logger}
}

// AllRequstHandler is handler which get all reqests by the userID.
// User id is getted from context.
// Handler calls service method GetRequest.
// Method response with json representation of request or error, if any occurs.
func (rh *ReqHandler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := jwt.GetUserFromContext(ctx)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	reqs, err := rh.requestService.GetRequests(ctx, userID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, reqs)
}

// AddRequstHandler is handler which add an reqest.
// Method response with request id or error, if any occurs.
// User id is getted from context.
// File is getted like a part from multipartForm.
// Information about convetsion is took like a part of multopartForm.
// Handler calls service method AddRequest.
func (rh *ReqHandler) AddRequest(w http.ResponseWriter, r *http.Request) {
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

	reqID, err := rh.requestService.AddRequest(ctx, userID, file, header.Filename, sendInfo)

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

// GetRequstHandler is handler which get one reqest.
// Method response with the json representation of request id or error, if any occurs.
// User id is getted from context.
// Request id is getted from query.
// Handler calls service method GetRequest.
func (rh *ReqHandler) GetRequest(w http.ResponseWriter, r *http.Request) {
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

	request, err := rh.requestService.GetRequest(ctx, userID, reqID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, request)
}

// DeleteRequstHandler is handler which delete one reqest.
// Method response with the id of deleted request or error, if any occurs.
// User id is getted from context.
// Request id is getted from query.
// Handler calls service method DeleteRequest.
func (rh *ReqHandler) DeleteRequest(w http.ResponseWriter, r *http.Request) {
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

	err = rh.requestService.DeleteRequest(ctx, userID, reqID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, reqID)
}

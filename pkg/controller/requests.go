package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	im "github.com/Dyleme/image-coverter"
	"github.com/gorilla/mux"
)

func (c *Controller) AllRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := c.getUserFromContext(r)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	reqs, err := c.service.GetRequests(userID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, reqs)
}

func (c *Controller) AddRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := c.getUserFromContext(r)
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

	var sendInfo im.ConversionInfo

	err = json.Unmarshal([]byte(info), &sendInfo)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var buf []byte

	_, err = file.Read(buf)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	imageID, err := c.service.AddRequest(userID, file, header.Filename, sendInfo)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	m := struct {
		ImageID int `json:"imageID"`
	}{
		ImageID: imageID,
	}
	newJSONResponse(w, m)
}

func (c *Controller) GetRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := c.getUserFromContext(r)
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

	request, err := c.service.GetRequest(userID, reqID)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newJSONResponse(w, request)
}

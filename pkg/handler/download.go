package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) DownloadImageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserFromContext(r)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)
	strImageID, ok := vars["id"]

	if !ok {
		newErrorResponse(w, http.StatusBadRequest, "id parameter is missing")
		return
	}

	imageID, err := strconv.Atoi(strImageID)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	b, err := s.service.DownloadImage(userID, imageID)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newDownloadFileResponce(w, b)
}

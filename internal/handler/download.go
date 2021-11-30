package handler

import (
	"net/http"
	"strconv"

	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/gorilla/mux"
)

func (h *Handler) DownloadImageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := jwt.GetUserFromContext(ctx)

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

	b, err := h.downloadService.DownloadImage(ctx, userID, imageID)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newDownloadFileResponce(w, b)
}
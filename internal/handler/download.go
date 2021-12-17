package handler

import (
	"net/http"
	"strconv"

	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Struct which provides method to handle downloading.
type DownHandler struct {
	logger          *logrus.Logger
	downloadService Downloader
}

// Constructor for DownHandler
func NewDownHandler(down Downloader, logger *logrus.Logger) *DownHandler {
	return &DownHandler{downloadService: down, logger: logger}
}

// DownloadImageHandler is Handler which response with image bytes.
// Handler get image id from query.
// Calls service method DownloadImage with image id and user id which i getted from context.
// If any error occurs than it response with error body.
func (dh *DownHandler) DownloadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := jwt.GetUserFromContext(ctx)
	if err != nil {
		newErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)
	strImageID, ok := vars["id"]

	if !ok {
		newErrorResponse(w, http.StatusBadRequest, `parameter "id" is missing`)
		return
	}

	imageID, err := strconv.Atoi(strImageID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	b, err := dh.downloadService.DownloadImage(ctx, userID, imageID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newDownloadFileResponse(w, b)
}

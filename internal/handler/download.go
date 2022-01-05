package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Downloader is an interface which has method to download image.
type Downloader interface {
	DownloadImage(ctx context.Context, userID, imageID int) ([]byte, string, error)
}

// Struct which provides method to handle downloading.
type Download struct {
	logger          *logrus.Logger
	downloadService Downloader
}

// Constructor for DownHandler.
func NewDownload(down Downloader, logger *logrus.Logger) *Download {
	return &Download{downloadService: down, logger: logger}
}

// DownloadImageHandler is Handler which response with image bytes.
// Handler get image id from query.
// Calls service method DownloadImage with image id and user id which i getted from context.
// If any error occurs than it response with error body.
func (dh *Download) DownloadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := jwt.GetUserFromContext(ctx)
	if err != nil {
		dh.logger.Warn(err)
		newErrorResponse(w, http.StatusUnauthorized, err.Error())

		return
	}

	vars := mux.Vars(r)
	strImageID, ok := vars["id"]

	if !ok {
		dh.logger.Warn(err)
		newErrorResponse(w, http.StatusBadRequest, `parameter "id" is missing`)

		return
	}

	imageID, err := strconv.Atoi(strImageID)
	if err != nil {
		dh.logger.Warn(err)
		newErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	b, filename, err := dh.downloadService.DownloadImage(ctx, userID, imageID)
	if err != nil {
		dh.logger.Warn(err)
		newErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	newDownloadFileResponse(w, b, filename)
}

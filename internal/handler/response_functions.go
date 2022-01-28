package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Struct to marshal error to the json.
type errorResponse struct {
	Message string `json:"message"`
}

func newUnknownErrorResponse(w http.ResponseWriter, err error) {
	var statErr ErrorWithStatus
	if errors.As(err, &statErr) {
		newErrorResponse(w, statErr.Status(), err.Error())
	} else {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// newErrorResponse functiton respnonse with the json representing the error.
func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	(w).Header().Set("Content-Type", "application/json; charset=utf-8")
	(w).Header().Set("X-Content-Type-Options", "nosniff")

	js, err := json.Marshal(errorResponse{message})

	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)

	fmt.Fprint(w, string(js))
}

// newJSONResponse function response with the json representing the interface v.
func newJSONResponse(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(js))
}

// newDownloadFileResponse response with bytes of files as attachment to the response.
func newDownloadFileResponse(w http.ResponseWriter, b []byte, filename string) {
	w.Header().Add("Content-Disposition", "Attachment")
	w.Header().Add("Content-Disposition", `filename="`+filename+`"`)
	w.Header().Add("Content-Length", strconv.Itoa(len(b)))
	fmt.Fprint(w, string(b))
}

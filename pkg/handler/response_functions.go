package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	log.Println(message)
	(w).Header().Set("Content-Type", "application/json; charset=utf-8")
	(w).Header().Set("X-Content-Type-Options", "nosniff")

	js, err := json.Marshal(errorResponse{message})

	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	fmt.Fprint(w, string(js))

	w.WriteHeader(statusCode)
}

func newJSONResponse(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(js))
}

func newDownloadFileResponce(w http.ResponseWriter, b []byte) {
	w.Header().Add("Content-Disposition", "Attachment")
	fmt.Fprint(w, string(b))
}

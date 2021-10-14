package handler

import "net/http"

func ConversionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("There should be a image conversion"))
}

package handler

import "net/http"

func HistortHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("There shoudld be aa History")) //nolint:errcheck // for future
}

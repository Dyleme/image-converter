package handler

import (
	"net/http"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELlo")) //nolint:errcheck // for future
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

}

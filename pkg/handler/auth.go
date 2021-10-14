package handler

import (
	"net/http"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELlo"))
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

}

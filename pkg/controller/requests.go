package controller

import "net/http"

func (c *Controller) AllRequestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("There shoudld be aa History")) //nolint:errcheck // for future
}

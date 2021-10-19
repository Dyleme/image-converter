package image

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,          // nolint:gomnd //one migabyte
		ReadTimeout:    10 * time.Second, // nolint:gomnd //10 seconds
		WriteTimeout:   10 * time.Second, // nolint:gomnd //10 seconds
	}

	return s.httpServer.ListenAndServe()
}

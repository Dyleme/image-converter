package image

import (
	"net/http"
	"time"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 10 * time.Second
	WriteTimeout   = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: MaxHeaderBytes,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
	}

	return s.httpServer.ListenAndServe()
}

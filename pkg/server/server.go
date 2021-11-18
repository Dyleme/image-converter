package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Dyleme/image-coverter/pkg/logging"
)

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second

	timeForGracefulShutdown = 5 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(ctx context.Context, port string, handler http.Handler) error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	logger := logging.FromContext(ctx)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		logger.Info("waiting system call")
		<-c
		logger.Info("system call")
		cancel()
	}()

	s.serve(ctx, port, handler)

	return nil
}

func (s *Server) serve(ctx context.Context, port string, handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	logger := logging.FromContext(ctx)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("listen %s", err)
		}
	}()

	logger.Info("server start")

	<-ctx.Done()

	logger.Info("server end")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), timeForGracefulShutdown)
	defer cancel()

	if err := s.httpServer.Shutdown(ctxShutDown); err != nil {
		logger.Fatalf("server shutdown failed %s", err)
	}

	logger.Info("server exited properly")
}

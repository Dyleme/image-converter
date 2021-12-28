package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Dyleme/image-coverter/internal/logging"
)

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second

	timeForGracefulShutdown = 5 * time.Second
)

// Server is a struct which handles the requests.
type Server struct {
	httpServer *http.Server
}

// After Run method Server starts to listen port and response to  the reqeusts.
// Run function provide the abitility of the gracefule shutdown.
func (s *Server) Run(ctx context.Context, port string, handler http.Handler) error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	logger := logging.FromContext(ctx)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-c
		logger.Info("system interrupt call")
		cancel()
	}()

	return s.serve(ctx, port, handler)
}

func (s *Server) serve(ctx context.Context, port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	logger := logging.FromContext(ctx)

	servError := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			servError <- fmt.Errorf("listen: %s", err)
		}
	}()

	logger.Info("server start")

	select {
	case err := <-servError:
		logger.Error("server crushed: ", err)
		return err

	case <-ctx.Done():
		logger.Info("server end")
		ctxShutDown, cancel := context.WithTimeout(context.Background(), timeForGracefulShutdown)
		defer cancel()

		if err := s.httpServer.Shutdown(ctxShutDown); err != nil {
			logger.Error("server didn't exit properly")
			return err
		}

		logger.Info("server exited properly")
	}

	return nil
}

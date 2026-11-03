package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

const (
	readHeaderTimeout = 10 * time.Second
	idleTimeout       = 90 * time.Second
	shutdownTimeout   = 15 * time.Second
)

type Server struct {
	Rtr *gin.Engine
}

func NewServer(env string, log *logrus.Logger) *Server {
	server := &Server{}
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	server.Rtr = gin.New()
	server.Rtr.Use(ginlogrus.Logger(log), gin.Recovery())
	server.registerMiddlewares()

	return server
}

func (s *Server) Start(ctx context.Context, port string) error {
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           s.Rtr,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		errCh <- err
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	// Останавливаем приём запросов и ждём завершения активных не дольше shutdownTimeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		_ = srv.Close()
		return fmt.Errorf("http server shutdown: %w", err)
	}
	return nil
}

package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/MyChaOS87/patAi/config"
	"github.com/MyChaOS87/patAi/internal/api/router"
	"github.com/MyChaOS87/patAi/pkg/log"
)

const (
	maxHeaderBytes = 1 << 20
)

type Server struct {
	api          *config.APIConfig
	echo         *echo.Echo
	childRouters []router.Router
}

func NewServer(options ...Option) *Server {
	cfg := newDefaultConfig()
	for _, opt := range options {
		opt(cfg)
	}

	return &Server{
		echo:         echo.New(),
		childRouters: cfg.childRouters,
		api:          cfg.api,
	}
}

// Run runs the server.
func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:              s.api.Server.Port,
		ReadTimeout:       s.api.Server.ReadTimeout,
		ReadHeaderTimeout: s.api.Server.ReadTimeout,
		WriteTimeout:      s.api.Server.WriteTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	if err := s.mapHandlers(); err != nil {
		return err
	}

	for _, r := range s.echo.Routes() {
		log.Infof("Route: %s", r.Path)
	}

	go func() {
		log.Infof("server is listening on PORT: %s", s.api.Server.Port)

		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("error starting Server: %v", err)
		}
	}()

	<-ctx.Done()

	srvCtx, shutdown := context.WithTimeout(context.Background(), s.api.Server.GracefulShutdownTimeout)
	defer shutdown()

	err := s.echo.Server.Shutdown(srvCtx)
	if err != nil {
		return errors.Wrap(err, "server graceful shutdown failed")
	}

	return nil
}

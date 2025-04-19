package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tehrelt/workmate-testovoe/internal/config"
	"github.com/tehrelt/workmate-testovoe/internal/transport/http/handlers"
	"github.com/tehrelt/workmate-testovoe/internal/transport/http/middlewares"
	"github.com/tehrelt/workmate-testovoe/pkg/sl"
)

type Server struct {
	router *echo.Echo
	cfg    *config.Config
}

func New(cfg *config.Config) *Server {
	s := &Server{
		router: echo.New(),
		cfg:    cfg,
	}

	return s.setup()
}

func (s *Server) setup() *Server {

	s.router.Use(middlewares.Tracing(s.cfg.Name))
	s.router.Use(middlewares.Logging)

	s.router.POST("/", handlers.CreateTask())

	return s
}

func (s *Server) Run(ctx context.Context) error {

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.router.Start(s.cfg.Http.Address()); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			slog.Error("failed to start http server", sl.Err(err))
			stop()
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()
	if err := s.router.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	slog.Info("server shutdown")

	return nil
}

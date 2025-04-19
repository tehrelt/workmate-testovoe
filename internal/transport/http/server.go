package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/tehrelt/workmate-testovoe/internal/config"
	"github.com/tehrelt/workmate-testovoe/internal/transport/http/handlers"
	"github.com/tehrelt/workmate-testovoe/internal/transport/http/middlewares"
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

	s.router.Use(middlewares.Logging)

	root := s.router.Group("/")
	root.POST("/", handlers.CreateTask())

	return s
}

func (s *Server) Run() error {
	if err := s.router.Start(s.cfg.Http.Address()); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

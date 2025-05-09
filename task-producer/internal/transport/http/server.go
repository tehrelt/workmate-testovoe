package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/config"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/transport/http/handlers"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/transport/http/middlewares"
	"github.com/tehrelt/workmate-testovoe/task-producer/pkg/sl"

	ut "github.com/go-playground/universal-translator"
)

type Server struct {
	router *echo.Echo
	cfg    *config.Config

	taskService *taskservice.TaskService
}

func New(cfg *config.Config, ts *taskservice.TaskService) *Server {
	s := &Server{
		router:      echo.New(),
		cfg:         cfg,
		taskService: ts,
	}

	return s.setup()
}

func (s *Server) setup() *Server {

	v := validator.New()
	uni := ut.New(en.New(), ru.New(), en.New())

	translator, _ := uni.GetTranslator("en")

	v.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	s.router.Validator = newHttpValidator(v, translator)

	s.router.Use(middlewares.Tracing(s.cfg.Name))
	s.router.Use(middlewares.Logging)
	s.router.Use(echoprometheus.NewMiddleware(s.cfg.Name))

	s.router.POST("/", handlers.CreateTask(s.taskService))
	s.router.GET("/", handlers.ListTasks(s.taskService))
	s.router.GET("/:id", handlers.GetTask(s.taskService))
	s.router.GET("/metrics", echoprometheus.NewHandler())

	return s
}

func (s *Server) Run(ctx context.Context) error {

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

		defer cancel()
		if err := s.router.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown HTTP server", sl.Err(err))
			return
		}

		slog.Info("server shutdown")
	}()

	if err := s.router.Start(s.cfg.Http.Address()); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}

		slog.Error("failed to start http server", sl.Err(err))
		return err
	}

	return nil
}

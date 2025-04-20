package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tehrelt/workmate-testovoe/task-producer/internal/config"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/transport/amqp"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/transport/http"
	"github.com/tehrelt/workmate-testovoe/task-producer/pkg/sl"
	"go.opentelemetry.io/otel/trace"
)

type Server interface {
	Run(context.Context) error
}

type App struct {
	cfg      *config.Config
	server   *http.Server
	consumer *amqp.Consumer
	tracer   trace.Tracer
}

func build(cfg *config.Config, server *http.Server, tracer trace.Tracer, consumer *amqp.Consumer) *App {
	return &App{
		cfg:      cfg,
		server:   server,
		tracer:   tracer,
		consumer: consumer,
	}
}

func (a *App) Run(ctx context.Context) error {

	servers := []Server{a.server, a.consumer}

	wg := sync.WaitGroup{}
	wg.Add(len(servers))

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, server := range servers {
		go func() {
			defer wg.Done()
			if err := server.Run(ctx); err != nil {
				slog.Error("failed to run server", sl.Err(err))
				stop()
			}
		}()
	}

	<-ctx.Done()
	wg.Wait()

	return nil
}

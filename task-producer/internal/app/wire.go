//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/config"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tracer"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage/pg/eventstorage"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage/pg/taskstorage"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/storage/rmq/taskqueue"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/transport/amqp"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/transport/http"
	"github.com/tehrelt/workmate-testovoe/task-producer/pkg/sl"
	"go.opentelemetry.io/otel/trace"
)

//go:generate wire .
func New(ctx context.Context) (*App, func(), error) {
	panic(
		wire.Build(
			build,

			http.New,
			_consumer,

			taskservice.New,
			wire.Bind(new(taskservice.TaskSaver), new(*taskstorage.TaskStorage)),
			wire.Bind(new(taskservice.TaskProvider), new(*taskstorage.TaskStorage)),
			wire.Bind(new(taskservice.TaskProcessor), new(*taskqueue.TaskQueue)),
			wire.Bind(new(taskservice.EventSaver), new(*eventstorage.Storage)),

			eventstorage.New,
			taskstorage.New,
			_newTasksProducer,

			_amqp,
			_pg,
			_tracer,
			config.New,
		),
	)
}

func _consumer(cfg *config.Config, ch *amqp091.Channel, ts *taskservice.TaskService) *amqp.Consumer {
	return amqp.New(ch, cfg.Queues.ProcessedTasks, ts)
}

func _newTasksProducer(cfg *config.Config, ch *amqp091.Channel) *taskqueue.TaskQueue {
	return taskqueue.New(cfg.Queues.NewTasks, ch)
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	jaeger := cfg.JaegerEndpoint
	appname := cfg.Name

	slog.Debug("connecting to jaeger", slog.String("jaeger", jaeger), slog.String("appname", appname))

	return tracer.SetupTracer(ctx, jaeger, appname)
}

func _pg(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, func(), error) {
	pool, err := pgxpool.Connect(ctx, cfg.PG.ConnectionString())
	if err != nil {
		return nil, nil, err
	}

	log := slog.With(slog.String("connection", cfg.PG.ConnectionString()))
	log.Debug("connecting to database")
	t := time.Now()
	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		return nil, func() { pool.Close() }, err
	}
	log.Info("connected to database", slog.String("ping", fmt.Sprintf("%2.fs", time.Since(t).Seconds())))

	return pool, func() { pool.Close() }, nil
}

func _amqp(cfg *config.Config) (*amqp091.Channel, func(), error) {

	slog.Debug("connecting to amqp", slog.String("cs", cfg.Amqp.ConnectionString()))
	conn, err := amqp091.Dial(cfg.Amqp.ConnectionString())
	if err != nil {
		slog.Error("failed to connect to amqp", sl.Err(err))
		return nil, nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		slog.Error("failed to create amqp channel", sl.Err(err))
		return nil, func() {
			conn.Close()
		}, err
	}

	closefn := func() {
		defer conn.Close()
		defer channel.Close()
	}

	if err := amqp_setup_exchange(channel, cfg.Queues.NewTasks.Exchange, cfg.Queues.NewTasks.RoutingKey); err != nil {
		slog.Error("failed to setup notifications exchange", sl.Err(err))
		return nil, closefn, err
	}
	if err := amqp_setup_exchange(channel, cfg.Queues.ProcessedTasks.Exchange, cfg.Queues.ProcessedTasks.RoutingKey); err != nil {
		slog.Error("failed to setup notifications exchange", sl.Err(err))
		return nil, closefn, err
	}

	return channel, closefn, nil
}

func amqp_setup_exchange(channel *amqp091.Channel, exchange string, queues ...string) error {

	log := slog.With(slog.String("exchange", exchange))
	log.Info("declaring exchange")
	if err := channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		slog.Error("failed to declare notifications queue", sl.Err(err))
		return err
	}

	for _, queueName := range queues {
		log.Info("declaring queue", slog.String("queue", queueName))
		queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			log.Error("failed to declare queue", sl.Err(err), slog.String("queue", queueName))
			return err
		}

		log.Info("binding queue", slog.String("queue", queueName))
		if err := channel.QueueBind(queue.Name, queueName, exchange, false, nil); err != nil {
			log.Error("failed to bind queue", sl.Err(err), slog.String("queue", queueName))
			return err
		}
	}

	return nil
}

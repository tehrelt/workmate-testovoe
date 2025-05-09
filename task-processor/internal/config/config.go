package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tehrelt/workmate-testovoe/task-processor/pkg/sl"
)

type Env string

const (
	EnvProd  Env = "prod"
	EnvDev   Env = "dev"
	EnvLocal Env = "local"
)

type Config struct {
	Env            Env        `env:"ENV" env-default:"local"`
	Name           string     `env:"APP_NAME" env-default:"workmate-testovoe"`
	Version        string     `env:"VERSION" env-default:"v0.1.0"`
	JaegerEndpoint string     `env:"JAEGER_ENDPOINT" env-default:"localhost:6831"`
	Amqp           AmqpConfig `env-prefix:"AMQP_" env-default:"localhost:5672:guest:guest:/"`
	PG             Database   `env-prefix:"PG_"`
	Queues         struct {
		NewTasks       QueueConfig `env-prefix:"NEW_TASKS_"`
		ProcessedTasks QueueConfig `env-prefix:"PROCESSED_TASKS_"`
	} `env-prefix:"QUEUE_"`
}

type AmqpConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	Port     int    `env:"PORT" env-default:"5672"`
	User     string `env:"USER" env-default:"guest"`
	Password string `env:"PASS" env-default:"guest"`
	Vhost    string `env:"VHOST" env-default:"/"`
}

func (a *AmqpConfig) ConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", a.User, a.Password, a.Host, a.Port)
}

type WorkerPoolConfig struct {
	MaxWorkers int `env:"MAX_WORKERS" env-default:"10"`
}

type Database struct {
	Protocol string `env:"PROTOCOL" env-default:"postgresql"`
	Host     string `env:"HOST" env-default:"localhost"`
	Port     int    `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-default:"postgres"`
	Password string `env:"PASS" env-default:"password"`
	Name     string `env:"NAME" env-default:"workmate"`
}

func (d *Database) ConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		d.Protocol, d.User, d.Password, d.Host, d.Port, d.Name)
}

type QueueConfig struct {
	Exchange         string           `env:"EXCHANGE" env-default:"workmate"`
	RoutingKey       string           `env:"ROUTING_KEY" env-default:"workmate"`
	WorkerPoolConfig WorkerPoolConfig `env-prefix:"POOL_" env-default:""`
}

func setupLogger(cfg *Config) {
	var log *slog.Logger

	switch cfg.Env {
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	slog.SetDefault(log)
}

func New() *Config {
	config := new(Config)

	if err := cleanenv.ReadEnv(config); err != nil {
		slog.Error("error when reading env", sl.Err(err))
		header := fmt.Sprintf("%s - %s", os.Getenv("APP_NAME"), os.Getenv("VERSION"))
		usage := cleanenv.FUsage(os.Stdout, config, &header)
		usage()

		os.Exit(-1)
	}

	setupLogger(config)

	slog.Debug("config", slog.Any("c", config))
	return config
}

package rmq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/workmate-testovoe/task-producer/internal/lib/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

type ConsumeFn func(ctx context.Context, msg amqp091.Delivery) error

type Manager struct {
	ch *amqp091.Channel
}

func New(ch *amqp091.Channel) *Manager {
	return &Manager{
		ch: ch,
	}
}

func (r *Manager) Close() error {
	return r.ch.Close()
}

func (r *Manager) Consume(ctx context.Context, rk string, fn ConsumeFn) error {
	messages, err := r.ch.ConsumeWithContext(ctx, rk, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range messages {
		carrier := propagation.MapCarrier{}
		for key, value := range msg.Headers {
			if v, ok := value.(string); ok {
				carrier.Set(key, v)
			}
		}
		p := propagation.TraceContext{}
		ctx = p.Extract(ctx, carrier)

		t := otel.Tracer(tracer.TracerKey)
		ctx, span := t.Start(ctx, fmt.Sprintf("Consume %s", rk))

		slog.Debug("consuming message", slog.Any("headers", msg.Headers))
		if err := fn(ctx, msg); err != nil {
			return err
		}

		span.End()
	}

	return nil
}

func (r *Manager) Publish(ctx context.Context, exchange, rk string, msg []byte) error {

	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, fmt.Sprintf("Publish %s.%s", exchange, rk))
	defer span.End()

	span.SetAttributes(
		attribute.String("amqp.exchange", exchange),
		attribute.String("amqp.routing_key", rk),
	)

	headers := make(amqp091.Table)

	propagator := propagation.TraceContext{}
	carrier := propagation.MapCarrier{}
	propagator.Inject(ctx, &carrier)

	for key, value := range carrier {
		headers[key] = value
	}

	slog.Debug("publishing message", slog.Any("headers", headers))

	return r.ch.PublishWithContext(ctx, exchange, rk, false, false, amqp091.Publishing{
		Headers:     headers,
		ContentType: "application/json",
		Body:        msg,
	})
}

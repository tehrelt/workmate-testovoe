package tx

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tehrelt/workmate-testovoe/task-processor/internal/lib/tracer"
	"go.opentelemetry.io/otel"
)

const (
	TxKey = "transaction"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrTransactionNotOpen  = errors.New("transaction not open")
)

type db interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type Transaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}

type TxOptions func(t *transaction) *transaction

func WithSpanName(name string) TxOptions {
	return func(t *transaction) *transaction {
		t.spanName = name
		return t
	}
}

func Begin(ctx context.Context, opts ...TxOptions) (context.Context, Transaction) {
	tracer := otel.Tracer(tracer.TracerKey)
	tx := newTx()

	for _, opt := range opts {
		tx = opt(tx)
	}

	ctx, span := tracer.Start(ctx, tx.spanName)
	tx.span = span
	return context.WithValue(ctx, TxKey, tx), tx
}

func get(ctx context.Context) *transaction {
	tx := ctx.Value(TxKey)
	if tx != nil {
		return tx.(*transaction)
	}

	return nil
}

func GetOrDefault(ctx context.Context, pool *pgxpool.Pool) (context.Context, db, error) {

	fn := "tx.GetOrDefault"
	log := slog.With(slog.String("fn", fn))

	tx := get(ctx)
	if tx == nil {
		log.Debug("tx not found in context, returning pool")
		return ctx, pool, nil
	}

	if tx.Tx != nil {
		log.Debug("tx found in context and pgx.Tx isnt nil, returning tx")
		return ctx, tx, nil
	}

	log.Debug("tx found in context, but pgx.Tx is nil, creating new pgx.Tx")
	pgtx, err := pool.Begin(ctx)
	if err != nil {
		return ctx, nil, err
	}

	log.Debug("enrich tx with pgx.Tx")
	tx.enrich(pgtx)

	return context.WithValue(ctx, TxKey, tx), tx, nil
}

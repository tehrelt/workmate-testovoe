package tx

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"go.opentelemetry.io/otel/trace"
)

type transaction struct {
	pgx.Tx
	span     trace.Span
	spanName string
}

func (t *transaction) Commit(ctx context.Context) error {
	defer t.span.End()
	return t.Tx.Commit(ctx)
}

func (t *transaction) Rollback(ctx context.Context) error {
	defer t.span.End()
	t.span.RecordError(errors.New("rollback tx"))
	return t.Tx.Rollback(ctx)
}

func (t *transaction) enrich(tx pgx.Tx) {
	t.Tx = tx
}

func newTx() *transaction {
	return &transaction{
		span:     nil,
		Tx:       nil,
		spanName: "transaction",
	}
}

package metrics

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability/tracing"

	"go.opentelemetry.io/otel/metric"
)

var _ UnitCounter = (*unitCounter)(nil)

type unitCounter struct {
	counter metric.Int64Counter
}

// NewUnitCounter yields a new UnitCounter.
func NewUnitCounter(counter metric.Int64Counter) UnitCounter {
	return &unitCounter{counter: counter}
}

func (c *unitCounter) Increment(ctx context.Context) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	c.counter.Add(ctx, 1)
}

func (c *unitCounter) IncrementBy(ctx context.Context, val int64) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	c.counter.Add(ctx, val)
}

func (c *unitCounter) Decrement(ctx context.Context) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	c.counter.Add(ctx, -1)
}

package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

type (
	Provider interface {
		NewFloat64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error)
		NewFloat64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error)
		NewFloat64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error)
		NewFloat64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error)
		NewInt64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error)
		NewInt64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error)
		NewInt64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error)
		NewInt64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error)
		Shutdown(ctx context.Context) error
	}
)

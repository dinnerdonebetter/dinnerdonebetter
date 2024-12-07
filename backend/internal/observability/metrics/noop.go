package metrics

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func NewNoopMetricsProvider() Provider {
	return &noopProvider{}
}

type noopProvider struct {
	mock.Mock
}

// NewFloat64Counter satisfies our interface
func (m *noopProvider) NewFloat64Counter(_ string, _ ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	return otel.Meter("noop").Float64Counter("")
}

// NewFloat64Gauge satisfies our interface
func (m *noopProvider) NewFloat64Gauge(_ string, _ ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	return otel.Meter("noop").Float64Gauge("")
}

// NewFloat64UpDownCounter satisfies our interface
func (m *noopProvider) NewFloat64UpDownCounter(_ string, _ ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	return otel.Meter("noop").Float64UpDownCounter("")
}

// NewFloat64Histogram satisfies our interface
func (m *noopProvider) NewFloat64Histogram(_ string, _ ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return otel.Meter("noop").Float64Histogram("")
}

// NewInt64Counter satisfies our interface
func (m *noopProvider) NewInt64Counter(_ string, _ ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return otel.Meter("noop").Int64Counter("")
}

// NewInt64Gauge satisfies our interface
func (m *noopProvider) NewInt64Gauge(_ string, _ ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return otel.Meter("noop").Int64Gauge("")
}

// NewInt64UpDownCounter satisfies our interface
func (m *noopProvider) NewInt64UpDownCounter(_ string, _ ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return otel.Meter("noop").Int64UpDownCounter("")
}

// NewInt64Histogram satisfies our interface
func (m *noopProvider) NewInt64Histogram(_ string, _ ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return otel.Meter("noop").Int64Histogram("")
}

// Shutdown satisfies our interface
func (m *noopProvider) Shutdown(ctx context.Context) error {
	return nil
}

package metrics

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/metric"
)

func NewNoopMetricsProvider() Provider {
	return &noopProvider{}
}

type noopProvider struct {
	mock.Mock
}

// NewFloat64Counter is a no-op method.
func (m *noopProvider) NewFloat64Counter(string, ...metric.Float64CounterOption) (Float64Counter, error) {
	return nil, nil
}

// NewFloat64Gauge is a no-op method.
func (m *noopProvider) NewFloat64Gauge(string, ...metric.Float64GaugeOption) (Float64Gauge, error) {
	return nil, nil
}

// NewFloat64UpDownCounter is a no-op method.
func (m *noopProvider) NewFloat64UpDownCounter(string, ...metric.Float64UpDownCounterOption) (Float64UpDownCounter, error) {
	return nil, nil
}

// NewFloat64Histogram is a no-op method.
func (m *noopProvider) NewFloat64Histogram(string, ...metric.Float64HistogramOption) (Float64Histogram, error) {
	return nil, nil
}

// NewInt64Counter is a no-op method.
func (m *noopProvider) NewInt64Counter(string, ...metric.Int64CounterOption) (Int64Counter, error) {
	return nil, nil
}

// NewInt64Gauge is a no-op method.
func (m *noopProvider) NewInt64Gauge(string, ...metric.Int64GaugeOption) (Int64Gauge, error) {
	return nil, nil
}

// NewInt64UpDownCounter is a no-op method.
func (m *noopProvider) NewInt64UpDownCounter(string, ...metric.Int64UpDownCounterOption) (Int64UpDownCounter, error) {
	return nil, nil
}

// NewInt64Histogram is a no-op method.
func (m *noopProvider) NewInt64Histogram(string, ...metric.Int64HistogramOption) (Int64Histogram, error) {
	return nil, nil
}

func (m *noopProvider) MeterProvider() metric.MeterProvider {
	return nil
}

// Shutdown satisfies our interface
func (m *noopProvider) Shutdown(ctx context.Context) error {
	return nil
}

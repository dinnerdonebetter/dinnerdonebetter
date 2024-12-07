package metrics

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/metric"
)

type mockProvider struct {
	mock.Mock
}

// NewFloat64Counter satisfies our interface
func (m *mockProvider) NewFloat64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	returnValues := m.Called(name, options)
	return returnValues.Get(0).(metric.Float64Counter), returnValues.Error(1)
}

// NewFloat64Gauge satisfies our interface
func (m *mockProvider) NewFloat64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	returnValues := m.Called(name, options)
	return returnValues.Get(0).(metric.Float64Gauge), returnValues.Error(1)
}

// NewFloat64UpDownCounter satisfies our interface
func (m *mockProvider) NewFloat64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	returnValues := m.Called(name, options)
	return returnValues.Get(0).(metric.Float64UpDownCounter), returnValues.Error(1)
}

// NewFloat64Histogram satisfies our interface
func (m *mockProvider) NewFloat64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	returnValues := m.Called(name, options)
	return returnValues.Get(0).(metric.Float64Histogram), returnValues.Error(1)
}

// NewInt64Counter satisfies our interface
func (m *mockProvider) NewInt64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	returnValues := m.Called(name, options)
	return returnValues.Get(0).(metric.Int64Counter), returnValues.Error(1)
}

// NewInt64Gauge satisfies our interface
func (m *mockProvider) NewInt64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	returnValues := m.Called(name, options)
	return returnValues.Get(0).(metric.Int64Gauge), returnValues.Error(1)
}

// NewInt64UpDownCounter satisfies our interface
func (m *mockProvider) NewInt64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	returnValues := m.Called(name, options)
	return returnValues.Get(0).(metric.Int64UpDownCounter), returnValues.Error(1)
}

// NewInt64Histogram satisfies our interface
func (m *mockProvider) NewInt64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	returnValues := m.Called(name, options)
	return returnValues.Get(0).(metric.Int64Histogram), returnValues.Error(1)
}

// Shutdown satisfies our interface
func (m *mockProvider) Shutdown(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

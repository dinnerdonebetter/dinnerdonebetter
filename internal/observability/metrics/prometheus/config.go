package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

const (
	// defaultNamespace is the default namespace under which we register metrics.
	defaultNamespace = "prixfixe_server"

	instrumentationVersion = "1.0.0"

	// minimumRuntimeCollectionInterval is the smallest interval we can collect metrics at
	// this value is used to guard against zero values.
	minimumRuntimeCollectionInterval = time.Second
)

type (
	// Config contains settings related to Prometheus.
	Config struct {
		_ struct{}

		// RuntimeMetricsCollectionInterval  is the interval we collect runtime statistics at.
		RuntimeMetricsCollectionInterval time.Duration `json:"runtimeMetricsCollectionInterval,omitempty" mapstructure:"runtime_metrics_collection_interval" toml:"runtime_metrics_collection_interval,omitempty"`
	}
)

// ValidateWithContext validates the config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.RuntimeMetricsCollectionInterval, validation.Min(minimumRuntimeCollectionInterval)),
	)
}

func (cfg *Config) initiateExporter() (*prometheus.Exporter, metric.MeterProvider, error) {
	config := prometheus.Config{
		// copied from go.opentelemetry.io/otel/sdk/metric/aggregator
		DefaultHistogramBoundaries: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)

	prometheusExporter, err := prometheus.New(config, c)
	if err != nil {
		return nil, nil, err
	}

	mp := prometheusExporter.MeterProvider()
	global.SetMeterProvider(mp)
	if err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(cfg.RuntimeMetricsCollectionInterval)); err != nil {
		return nil, nil, fmt.Errorf("failed to start runtime metrics collection: %w", err)
	}

	return prometheusExporter, mp, nil
}

// ProvideMetricsHandler provides an instrumentation handler.
func (cfg *Config) ProvideMetricsHandler() (metrics.Handler, error) {
	mh, _, err := cfg.initiateExporter()

	return mh, err
}

// ProvideUnitCounterProvider provides an instrumentation handler.
func (cfg *Config) ProvideUnitCounterProvider(logger logging.Logger) (metrics.UnitCounterProvider, error) {
	logger.Debug("setting up meter")

	_, meterProvider, err := cfg.initiateExporter()
	if err != nil {
		return nil, err
	}

	mustMeter := metric.Must(meterProvider.Meter(defaultNamespace, metric.WithInstrumentationVersion(instrumentationVersion)))
	logger.Debug("meter initialized successfully")

	uc := func(name, description string) metrics.UnitCounter {
		l := logger.WithValue("name", name)

		counter := mustMeter.NewInt64Counter(
			fmt.Sprintf("%s_count", name),
			metric.WithUnit(unit.Dimensionless),
			metric.WithDescription(description),
		)

		l.Debug("returning wrapped unit counter")

		return metrics.NewUnitCounter(counter)
	}

	return uc, nil
}

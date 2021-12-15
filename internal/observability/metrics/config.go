package metrics

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

const (
	// defaultNamespace is the default namespace under which we register metrics.
	defaultNamespace = "prixfixe_server"

	instrumentationVersion = "1.0.0"

	// minimumRuntimeCollectionInterval is the smallest interval we can collect metrics at
	// this value is used to guard against zero values.
	minimumRuntimeCollectionInterval = time.Second

	// DefaultMetricsCollectionInterval is the default amount of time we wait between runtime metrics queries.
	DefaultMetricsCollectionInterval = 2 * time.Second

	// Prometheus represents the popular time series database.
	Prometheus = "prometheus"
)

type (
	// Config contains settings related to .
	Config struct {
		_ struct{}

		// Provider indicates where our metrics should go.
		Provider string `json:"provider,omitempty" mapstructure:"provider" toml:"provider,omitempty"`
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

var (
	// regretful state artifacts.
	prometheusExporterInitOnce sync.Once
	prometheusExporter         *otelprom.Exporter
)

func initiatePrometheusExporter() {
	prometheusExporterInitOnce.Do(func() {
		var err error

		config := otelprom.Config{
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

		if prometheusExporter, err = otelprom.New(config, c); err != nil {
			panic(err)
		}
	})
}

// ProvideInstrumentationHandler provides an instrumentation handler.
func (cfg *Config) ProvideInstrumentationHandler(logger logging.Logger) (InstrumentationHandler, error) {
	logger = logger.WithValue("metrics_provider", cfg.Provider)
	logger.Debug("setting metrics provider")

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(cfg.RuntimeMetricsCollectionInterval)); err != nil {
		return nil, fmt.Errorf("failed to start runtime metrics collection: %w", err)
	}

	switch strings.TrimSpace(strings.ToLower(cfg.Provider)) {
	case Prometheus:
		initiatePrometheusExporter()
		return prometheusExporter, nil
	default:
		// not a crime to not need metrics
		return nil, nil
	}
}

// ProvideUnitCounterProvider provides an instrumentation handler.
func (cfg *Config) ProvideUnitCounterProvider(logger logging.Logger) (UnitCounterProvider, error) {
	p := strings.TrimSpace(strings.ToLower(cfg.Provider))

	logger = logger.WithValue("metrics_provider", p)
	logger.Debug("setting up meter")

	switch p {
	case Prometheus:
		initiatePrometheusExporter()

		meterProvider := prometheusExporter.MeterProvider()
		mustMeter := metric.Must(meterProvider.Meter(defaultNamespace, metric.WithInstrumentationVersion(instrumentationVersion)))

		logger.Debug("meter initialized successfully")

		return func(name, description string) UnitCounter {
			l := logger.WithValue("name", name)

			counter := mustMeter.NewInt64Counter(
				fmt.Sprintf("%s_count", name),
				metric.WithUnit(unit.Dimensionless),
				metric.WithDescription(description),
			)

			l.Debug("returning wrapped unit counter")

			return &unitCounter{counter: counter}
		}, nil
	default:
		logger.Debug("nil unit counter provider")
		// not a crime to not need metrics
		return nil, nil
	}
}

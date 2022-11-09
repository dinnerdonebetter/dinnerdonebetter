package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/metrics"
	"github.com/prixfixeco/backend/internal/observability/metrics/prometheus"
)

const (
	// DefaultMetricsCollectionInterval is the default amount of time we wait between runtime metrics queries.
	DefaultMetricsCollectionInterval = 2 * time.Second

	// ProviderPrometheus represents the popular time series database.
	ProviderPrometheus = "prometheus"
)

type (
	// Config contains settings related to .
	Config struct {
		_          struct{}
		Prometheus *prometheus.Config `json:"prometheus,omitempty" mapstructure:"prometheus" toml:"prometheus,omitempty"`
		Provider   string             `json:"provider,omitempty" mapstructure:"provider" toml:"provider,omitempty"`
	}
)

// ValidateWithContext validates the config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In("", ProviderPrometheus)),
		validation.Field(&cfg.Prometheus, validation.When(cfg.Provider == ProviderPrometheus, validation.Required)),
	)
}

// ProvideMetricsHandler provides an instrumentation handler.
func (cfg *Config) ProvideMetricsHandler(l logging.Logger) (metrics.Handler, error) {
	p := strings.TrimSpace(strings.ToLower(cfg.Provider))

	logger := l.WithValue("metrics_provider", p)
	logger.Debug("setting up meter handler")

	switch p {
	case ProviderPrometheus:
		return cfg.Prometheus.ProvideMetricsHandler()
	default:
		return nil, nil
	}
}

// ProvideUnitCounterProvider provides a counter provider.
func (cfg *Config) ProvideUnitCounterProvider(_ context.Context, logger logging.Logger) (metrics.UnitCounterProvider, error) {
	p := strings.TrimSpace(strings.ToLower(cfg.Provider))

	logger = logger.WithValue("metrics_provider", p)
	logger.Debug("setting up meter provider")

	switch p {
	case ProviderPrometheus:
		return cfg.Prometheus.ProvideUnitCounterProvider(logger)
	case "":
		logger.Debug("noop unit counter provider")
		noopFunc := metrics.UnitCounterProvider(func(name, description string) metrics.UnitCounter {
			return metrics.EnsureUnitCounter(nil, logger, "", "")
		})

		return noopFunc, nil
	default:
		return nil, fmt.Errorf("invalid provider: %q", p)
	}
}

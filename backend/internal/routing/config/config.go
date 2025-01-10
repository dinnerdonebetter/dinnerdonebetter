package routingcfg

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderChi is the string we use to refer to chi.
	ProviderChi = "chi"
)

// Config configures our router.
type Config struct {
	_ struct{} `json:"-"`

	ChiConfig *chi.Config `envPrefix:"CHI_" json:"chiConfig,omitempty"`
	Provider  string      `env:"PROVIDER"   json:"provider,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a router config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In(ProviderChi)),
	)
}

// ProvideRouter provides a Router from a routing config.
func ProvideRouter(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, metricProvider metrics.Provider) (routing.Router, error) {
	switch cfg.Provider {
	case ProviderChi:
		return chi.NewRouter(logger, tracerProvider, metricProvider, cfg.ChiConfig), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}

// ProvideRouteParamManager provides a RouteParamManager from a routing config.
func ProvideRouteParamManager(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, metricProvider metrics.Provider) (routing.RouteParamManager, error) {
	switch cfg.Provider {
	case ProviderChi:
		return chi.NewRouteParamManager(), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}

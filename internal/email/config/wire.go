package config

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/google/wire"
)

var (
	// ProvidersEmail are what we provide to dependency injection.
	ProvidersEmail = wire.NewSet(
		ProvideEmailer,
	)
)

// ProvideEmailer provides an email.Emailer from a config.
func ProvideEmailer(ctx context.Context, cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client) (email.Emailer, error) {
	return cfg.ProvideEmailer(ctx, logger, tracerProvider, client)
}

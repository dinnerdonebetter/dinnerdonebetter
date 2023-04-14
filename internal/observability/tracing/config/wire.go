package config

import (
	"context"

	"github.com/google/wire"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

var (
	// Providers is a Wire provider set that provides a tracing.TracerProvider.
	Providers = wire.NewSet(
		ProvideTracerProvider,
	)
)

func ProvideTracerProvider(ctx context.Context, c *Config, l logging.Logger) (traceProvider tracing.TracerProvider, err error) {
	return c.ProvideTracerProvider(ctx, l)
}
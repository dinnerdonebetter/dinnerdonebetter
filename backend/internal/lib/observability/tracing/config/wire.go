package tracingcfg

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/google/wire"
)

var (
	// ProvidersTracingConfig is a Wire provider set that provides a tracing.TracerProvider.
	ProvidersTracingConfig = wire.NewSet(
		ProvideTracerProvider,
	)
)

func ProvideTracerProvider(ctx context.Context, c *Config, l logging.Logger) (traceProvider tracing.TracerProvider, err error) {
	return c.ProvideTracerProvider(ctx, l)
}

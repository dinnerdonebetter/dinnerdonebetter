package vectorcfg

import (
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/search/vector"

	"github.com/google/wire"
)

var (
	VectorProviders = wire.NewSet(
		ProvideVectorSearcher,
	)
)

func ProvideVectorSearcher(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider) (vector.Searcher, error) {
	return cfg.ProvideVectorSearcher(logger, tracerProvider)
}

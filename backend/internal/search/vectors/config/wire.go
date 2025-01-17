package vectorcfg

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search/vectors"

	"github.com/google/wire"
)

var (
	VectorProviders = wire.NewSet(
		ProvideVectorSearcher,
	)
)

func ProvideVectorSearcher(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider) (vectors.Searcher, error) {
	return cfg.ProvideVectorSearcher(logger, tracerProvider)
}

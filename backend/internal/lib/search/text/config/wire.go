package textsearchcfg

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideValidIngredientIndexManager,
	)
)

func ProvideValidIngredientIndexManager(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, metricsProvider metrics.Provider, cfg *Config) (textsearch.Index[types.ValidIngredientSearchSubset], error) {
	return ProvideIndex[types.ValidIngredientSearchSubset](ctx, logger, tracerProvider, metricsProvider, cfg, textsearch.IndexTypeValidIngredients)
}

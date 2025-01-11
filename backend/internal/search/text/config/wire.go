package textsearchcfg

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search/text"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideValidIngredientIndexManager,
	)
)

func ProvideValidIngredientIndexManager(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config) (textsearch.Index[types.ValidIngredientSearchSubset], error) {
	return ProvideIndex[types.ValidIngredientSearchSubset](ctx, logger, tracerProvider, cfg, textsearch.IndexTypeValidIngredients)
}

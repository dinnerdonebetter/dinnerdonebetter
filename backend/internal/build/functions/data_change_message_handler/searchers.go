package datachangemessagehandler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/core/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"

	"github.com/google/wire"
)

var SearcherProviders = wire.NewSet(
	ProvideUserTextSearcher,
	ProvideRecipeTextSearcher,
	ProvideMealTextSearcher,
	ProvideValidIngredientTextSearcher,
	ProvideValidInstrumentTextSearcher,
	ProvideValidMeasurementUnitTextSearcher,
	ProvideValidPreparationTextSearcher,
	ProvideValidIngredientStateTextSearcher,
	ProvideValidVesselTextSearcher,
)

func ProvideUserTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (coreindexing.UserTextSearcher, error) {
	return textsearchcfg.ProvideIndex[coreindexing.UserSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		coreindexing.IndexTypeUsers,
	)
}

func ProvideRecipeTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (eatingindexing.RecipeTextSearcher, error) {
	return textsearchcfg.ProvideIndex[eatingindexing.RecipeSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		eatingindexing.IndexTypeRecipes,
	)
}

func ProvideMealTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (eatingindexing.MealTextSearcher, error) {
	return textsearchcfg.ProvideIndex[eatingindexing.MealSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		eatingindexing.IndexTypeMeals,
	)
}

func ProvideValidIngredientTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (eatingindexing.ValidIngredientTextSearcher, error) {
	return textsearchcfg.ProvideIndex[eatingindexing.ValidIngredientSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		eatingindexing.IndexTypeValidIngredients,
	)
}

func ProvideValidInstrumentTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (eatingindexing.ValidInstrumentTextSearcher, error) {
	return textsearchcfg.ProvideIndex[eatingindexing.ValidInstrumentSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		eatingindexing.IndexTypeValidInstruments,
	)
}

func ProvideValidMeasurementUnitTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (eatingindexing.ValidMeasurementUnitTextSearcher, error) {
	return textsearchcfg.ProvideIndex[eatingindexing.ValidMeasurementUnitSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		eatingindexing.IndexTypeValidMeasurementUnits,
	)
}

func ProvideValidPreparationTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (eatingindexing.ValidPreparationTextSearcher, error) {
	return textsearchcfg.ProvideIndex[eatingindexing.ValidPreparationSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		eatingindexing.IndexTypeValidPreparations,
	)
}

func ProvideValidIngredientStateTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (eatingindexing.ValidIngredientStateTextSearcher, error) {
	return textsearchcfg.ProvideIndex[eatingindexing.ValidIngredientStateSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		eatingindexing.IndexTypeValidIngredientStates,
	)
}

func ProvideValidVesselTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (eatingindexing.ValidVesselTextSearcher, error) {
	return textsearchcfg.ProvideIndex[eatingindexing.ValidVesselSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		eatingindexing.IndexTypeValidVessels,
	)
}

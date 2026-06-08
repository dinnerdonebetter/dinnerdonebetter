package indexing

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"

	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/metrics"
	"github.com/primandproper/platform/observability/tracing"
	textsearchcfg "github.com/primandproper/platform/search/text/config"

	"github.com/samber/do/v2"
)

// RegisterMealPlanningSearchers registers all mealplanning text searcher providers with the injector.
func RegisterMealPlanningSearchers(i do.Injector) {
	do.Provide(i, func(i do.Injector) (RecipeTextSearcher, error) {
		return textsearchcfg.ProvideIndex[RecipeSearchSubset](
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			IndexTypeRecipes,
		)
	})
	do.Provide(i, func(i do.Injector) (MealTextSearcher, error) {
		return textsearchcfg.ProvideIndex[MealSearchSubset](
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			IndexTypeMeals,
		)
	})
	do.Provide(i, func(i do.Injector) (ValidIngredientTextSearcher, error) {
		return textsearchcfg.ProvideIndex[ValidIngredientSearchSubset](
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			IndexTypeValidIngredients,
		)
	})
	do.Provide(i, func(i do.Injector) (ValidInstrumentTextSearcher, error) {
		return textsearchcfg.ProvideIndex[ValidInstrumentSearchSubset](
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			IndexTypeValidInstruments,
		)
	})
	do.Provide(i, func(i do.Injector) (ValidMeasurementUnitTextSearcher, error) {
		return textsearchcfg.ProvideIndex[ValidMeasurementUnitSearchSubset](
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			IndexTypeValidMeasurementUnits,
		)
	})
	do.Provide(i, func(i do.Injector) (ValidPreparationTextSearcher, error) {
		return textsearchcfg.ProvideIndex[ValidPreparationSearchSubset](
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			IndexTypeValidPreparations,
		)
	})
	do.Provide(i, func(i do.Injector) (ValidIngredientStateTextSearcher, error) {
		return textsearchcfg.ProvideIndex[ValidIngredientStateSearchSubset](
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			IndexTypeValidIngredientStates,
		)
	})
	do.Provide(i, func(i do.Injector) (ValidVesselTextSearcher, error) {
		return textsearchcfg.ProvideIndex[ValidVesselSearchSubset](
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*textsearchcfg.Config](i),
			IndexTypeValidVessels,
		)
	})
}

// RegisterMealPlanningDataIndexer registers the meal planning data indexer with the injector.
func RegisterMealPlanningDataIndexer(i do.Injector) {
	do.Provide[*MealPlanningDataIndexer](i, func(i do.Injector) (*MealPlanningDataIndexer, error) {
		return NewMealPlanningDataIndexer(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[mealplanning.Repository](i),
			do.MustInvoke[RecipeTextSearcher](i),
			do.MustInvoke[MealTextSearcher](i),
			do.MustInvoke[ValidIngredientTextSearcher](i),
			do.MustInvoke[ValidInstrumentTextSearcher](i),
			do.MustInvoke[ValidMeasurementUnitTextSearcher](i),
			do.MustInvoke[ValidPreparationTextSearcher](i),
			do.MustInvoke[ValidIngredientStateTextSearcher](i),
			do.MustInvoke[ValidVesselTextSearcher](i),
		), nil
	})
}

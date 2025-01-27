package indexing

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	ErrNilIndexRequest = errors.New("nil index request")
)

func HandleIndexRequest(
	ctx context.Context,
	l logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	searchConfig *textsearchcfg.Config,
	dataManager database.DataManager,
	indexReq *textsearch.IndexRequest,
) error {
	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("eating_search_indexer"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	if indexReq == nil {
		return observability.PrepareAndLogError(ErrNilIndexRequest, l, span, "handling index requests")
	}

	logger := l.WithValue("index_type_requested", indexReq.IndexType)

	var (
		im                textsearch.IndexManager
		toBeIndexed       any
		err               error
		markAsIndexedFunc func() error
	)

	switch indexReq.IndexType {
	case textsearch.IndexTypeRecipes:
		im, err = textsearchcfg.ProvideIndex[RecipeSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var recipe *types.Recipe
		recipe, err = dataManager.GetRecipe(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting recipe")
		}

		toBeIndexed = ConvertRecipeToRecipeSearchSubset(recipe)
		markAsIndexedFunc = func() error { return dataManager.MarkRecipeAsIndexed(ctx, indexReq.RowID) }

	case textsearch.IndexTypeMeals:
		im, err = textsearchcfg.ProvideIndex[MealSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var meal *types.Meal
		meal, err = dataManager.GetMeal(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting meal")
		}

		toBeIndexed = ConvertMealToMealSearchSubset(meal)
		markAsIndexedFunc = func() error { return dataManager.MarkMealAsIndexed(ctx, indexReq.RowID) }

	case textsearch.IndexTypeValidIngredients:
		im, err = textsearchcfg.ProvideIndex[ValidIngredientSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validIngredient *types.ValidIngredient
		validIngredient, err = dataManager.GetValidIngredient(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient")
		}

		toBeIndexed = ConvertValidIngredientToValidIngredientSearchSubset(validIngredient)
		markAsIndexedFunc = func() error { return dataManager.MarkValidIngredientAsIndexed(ctx, indexReq.RowID) }

	case textsearch.IndexTypeValidInstruments:
		im, err = textsearchcfg.ProvideIndex[ValidInstrumentSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validInstrument *types.ValidInstrument
		validInstrument, err = dataManager.GetValidInstrument(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid instrument")
		}

		toBeIndexed = ConvertValidInstrumentToValidInstrumentSearchSubset(validInstrument)
		markAsIndexedFunc = func() error { return dataManager.MarkValidInstrumentAsIndexed(ctx, indexReq.RowID) }

	case textsearch.IndexTypeValidMeasurementUnits:
		im, err = textsearchcfg.ProvideIndex[ValidMeasurementUnitSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validMeasurementUnit *types.ValidMeasurementUnit
		validMeasurementUnit, err = dataManager.GetValidMeasurementUnit(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid measurement unit")
		}

		toBeIndexed = ConvertValidMeasurementUnitToValidMeasurementUnitSearchSubset(validMeasurementUnit)
		markAsIndexedFunc = func() error { return dataManager.MarkValidMeasurementUnitAsIndexed(ctx, indexReq.RowID) }

	case textsearch.IndexTypeValidPreparations:
		im, err = textsearchcfg.ProvideIndex[ValidPreparationSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validPreparation *types.ValidPreparation
		validPreparation, err = dataManager.GetValidPreparation(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid preparation")
		}

		toBeIndexed = ConvertValidPreparationToValidPreparationSearchSubset(validPreparation)
		markAsIndexedFunc = func() error { return dataManager.MarkValidPreparationAsIndexed(ctx, indexReq.RowID) }

	case textsearch.IndexTypeValidIngredientStates:
		im, err = textsearchcfg.ProvideIndex[ValidIngredientStateSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validIngredientState *types.ValidIngredientState
		validIngredientState, err = dataManager.GetValidIngredientState(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient state")
		}

		toBeIndexed = ConvertValidIngredientStateToValidIngredientStateSearchSubset(validIngredientState)
		markAsIndexedFunc = func() error { return dataManager.MarkValidIngredientStateAsIndexed(ctx, indexReq.RowID) }

	case textsearch.IndexTypeValidVessels:
		im, err = textsearchcfg.ProvideIndex[ValidVesselSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validVessel *types.ValidVessel
		validVessel, err = dataManager.GetValidVessel(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid vessel")
		}

		toBeIndexed = ConvertValidVesselToValidVesselSearchSubset(validVessel)
		markAsIndexedFunc = func() error { return dataManager.MarkValidVesselAsIndexed(ctx, indexReq.RowID) }
	default:
		logger.Info("invalid index type specified, exiting")
		return nil
	}

	if indexReq.Delete {
		if err = im.Delete(ctx, indexReq.RowID); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "deleting data")
		}

		return nil
	} else {
		if err = im.Index(ctx, indexReq.RowID, toBeIndexed); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "indexing data")
		}

		if err = markAsIndexedFunc(); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "marking data as indexed")
		}
	}

	return nil
}

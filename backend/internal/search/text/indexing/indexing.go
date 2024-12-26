package indexing

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

var (
	ErrNilIndexRequest = errors.New("nil index request")

	// AllIndexTypes is a list of all index types.
	AllIndexTypes = []string{
		textsearch.IndexTypeRecipes,
		textsearch.IndexTypeMeals,
		textsearch.IndexTypeValidIngredients,
		textsearch.IndexTypeValidInstruments,
		textsearch.IndexTypeValidMeasurementUnits,
		textsearch.IndexTypeValidPreparations,
		textsearch.IndexTypeValidIngredientStates,
		textsearch.IndexTypeValidVessels,
		textsearch.IndexTypeUsers,
	}
)

func HandleIndexRequest(ctx context.Context, l logging.Logger, tracerProvider tracing.TracerProvider, searchConfig *textsearchcfg.Config, dataManager database.DataManager, indexReq *IndexRequest) error {
	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("search-indexer"))
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
		im, err = textsearchcfg.ProvideIndex[types.RecipeSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var recipe *types.Recipe
		recipe, err = dataManager.GetRecipe(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting recipe")
		}

		toBeIndexed = converters.ConvertRecipeToRecipeSearchSubset(recipe)
		markAsIndexedFunc = func() error { return dataManager.MarkRecipeAsIndexed(ctx, indexReq.RowID) }
	case textsearch.IndexTypeMeals:
		im, err = textsearchcfg.ProvideIndex[types.MealSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var meal *types.Meal
		meal, err = dataManager.GetMeal(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting meal")
		}

		toBeIndexed = converters.ConvertMealToMealSearchSubset(meal)
		markAsIndexedFunc = func() error { return dataManager.MarkMealAsIndexed(ctx, indexReq.RowID) }
	case textsearch.IndexTypeValidIngredients:
		im, err = textsearchcfg.ProvideIndex[types.ValidIngredientSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validIngredient *types.ValidIngredient
		validIngredient, err = dataManager.GetValidIngredient(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient")
		}

		toBeIndexed = converters.ConvertValidIngredientToValidIngredientSearchSubset(validIngredient)
		markAsIndexedFunc = func() error { return dataManager.MarkValidIngredientAsIndexed(ctx, indexReq.RowID) }
	case textsearch.IndexTypeValidInstruments:
		im, err = textsearchcfg.ProvideIndex[types.ValidInstrumentSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validInstrument *types.ValidInstrument
		validInstrument, err = dataManager.GetValidInstrument(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid instrument")
		}

		toBeIndexed = converters.ConvertValidInstrumentToValidInstrumentSearchSubset(validInstrument)
		markAsIndexedFunc = func() error { return dataManager.MarkValidInstrumentAsIndexed(ctx, indexReq.RowID) }
	case textsearch.IndexTypeValidMeasurementUnits:
		im, err = textsearchcfg.ProvideIndex[types.ValidMeasurementUnitSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validMeasurementUnit *types.ValidMeasurementUnit
		validMeasurementUnit, err = dataManager.GetValidMeasurementUnit(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid measurement unit")
		}

		toBeIndexed = converters.ConvertValidMeasurementUnitToValidMeasurementUnitSearchSubset(validMeasurementUnit)
		markAsIndexedFunc = func() error { return dataManager.MarkValidMeasurementUnitAsIndexed(ctx, indexReq.RowID) }
	case textsearch.IndexTypeValidPreparations:
		im, err = textsearchcfg.ProvideIndex[types.ValidPreparationSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validPreparation *types.ValidPreparation
		validPreparation, err = dataManager.GetValidPreparation(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid preparation")
		}

		toBeIndexed = converters.ConvertValidPreparationToValidPreparationSearchSubset(validPreparation)
		markAsIndexedFunc = func() error { return dataManager.MarkValidPreparationAsIndexed(ctx, indexReq.RowID) }
	case textsearch.IndexTypeValidIngredientStates:
		im, err = textsearchcfg.ProvideIndex[types.ValidIngredientStateSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validIngredientState *types.ValidIngredientState
		validIngredientState, err = dataManager.GetValidIngredientState(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient state")
		}

		toBeIndexed = converters.ConvertValidIngredientStateToValidIngredientStateSearchSubset(validIngredientState)
		markAsIndexedFunc = func() error { return dataManager.MarkValidIngredientStateAsIndexed(ctx, indexReq.RowID) }
	case textsearch.IndexTypeValidVessels:
		im, err = textsearchcfg.ProvideIndex[types.ValidVesselSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validVessel *types.ValidVessel
		validVessel, err = dataManager.GetValidVessel(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid vessel")
		}

		toBeIndexed = converters.ConvertValidVesselToValidVesselSearchSubset(validVessel)
		markAsIndexedFunc = func() error { return dataManager.MarkValidVesselAsIndexed(ctx, indexReq.RowID) }
	case textsearch.IndexTypeUsers:
		im, err = textsearchcfg.ProvideIndex[types.UserSearchSubset](ctx, logger, tracerProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var user *types.User
		user, err = dataManager.GetUser(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting user")
		}

		toBeIndexed = converters.ConvertUserToUserSearchSubset(user)
		markAsIndexedFunc = func() error { return dataManager.MarkUserAsIndexed(ctx, indexReq.RowID) }
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

type IndexRequest struct {
	RequestID string `json:"id"`
	RowID     string `json:"rowID"`
	IndexType string `json:"type"`
	Delete    bool   `json:"delete"`
}

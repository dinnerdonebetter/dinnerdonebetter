package indexing

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	// IndexTypeRecipes represents the recipes index.
	IndexTypeRecipes = "recipes"
	// IndexTypeMeals represents the meals index.
	IndexTypeMeals = "meals"
	// IndexTypeValidIngredients represents the valid_ingredients index.
	IndexTypeValidIngredients = "valid_ingredients"
	// IndexTypeValidInstruments represents the valid_instruments index.
	IndexTypeValidInstruments = "valid_instruments"
	// IndexTypeValidMeasurementUnits represents the valid_measurement_units index.
	IndexTypeValidMeasurementUnits = "valid_measurement_units"
	// IndexTypeValidPreparations represents the  valid_preparations index.
	IndexTypeValidPreparations = "valid_preparations"
	// IndexTypeValidIngredientStates represents the valid_ingredient_states index.
	IndexTypeValidIngredientStates = "valid_ingredient_states"
	// IndexTypeValidIngredientMeasurementUnits represents the valid_ingredient_measurement_units index.
	IndexTypeValidIngredientMeasurementUnits = "valid_ingredient_measurement_units"
	// IndexTypeValidPreparationInstruments represents the valid_preparation_instruments index.
	IndexTypeValidPreparationInstruments = "valid_preparation_instruments"
	// IndexTypeValidIngredientPreparations represents the valid_ingredient_preparations index.
	IndexTypeValidIngredientPreparations = "valid_ingredient_preparations"
	// IndexTypeValidMeasurementUnitConversions represents the valid_measurement_unit_conversions index.
	IndexTypeValidMeasurementUnitConversions = "valid_measurement_unit_conversions"
)

var (
	ErrNilIndexRequest = errors.New("nil index request")
)

func HandleIndexRequest(ctx context.Context, l logging.Logger, tracerProvider tracing.TracerProvider, searchConfig *config.Config, dataManager database.DataManager, searchIndexRequest *IndexRequest) error {
	tracer := tracing.NewTracer(tracerProvider.Tracer("search-indexer"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	if searchIndexRequest == nil {
		return observability.PrepareAndLogError(ErrNilIndexRequest, l, span, "handling index requests")
	}

	logger := l.WithValue("index_type_requested", searchIndexRequest.IndexType)

	var (
		im          search.IndexManager
		toBeIndexed any
		err         error
	)

	switch searchIndexRequest.IndexType {
	case IndexTypeRecipes:
		im, err = config.ProvideIndexManager[types.RecipeSearchSubset](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var recipe *types.Recipe
		recipe, err = dataManager.GetRecipe(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting meal")
		}

		toBeIndexed = converters.ConvertRecipeToRecipeSearchSubset(recipe)
		if err = im.Index(ctx, recipe.ID, toBeIndexed); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "indexing meal")
		}

		return nil
	case IndexTypeMeals:
		im, err = config.ProvideIndexManager[types.MealSearchSubset](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var meal *types.Meal
		meal, err = dataManager.GetMeal(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting meal")
		}

		toBeIndexed = converters.ConvertMealToMealSearchSubset(meal)
	case IndexTypeValidIngredients:
		im, err = config.ProvideIndexManager[types.ValidIngredient](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validIngredient *types.ValidIngredient
		validIngredient, err = dataManager.GetValidIngredient(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient")
		}

		toBeIndexed = converters.ConvertValidIngredientToValidIngredientSearchSubset(validIngredient)
	case IndexTypeValidInstruments:
		im, err = config.ProvideIndexManager[types.ValidInstrument](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validInstrument *types.ValidInstrument
		validInstrument, err = dataManager.GetValidInstrument(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid instrument")
		}

		toBeIndexed = converters.ConvertValidInstrumentToValidInstrumentSearchSubset(validInstrument)
	case IndexTypeValidMeasurementUnits:
		im, err = config.ProvideIndexManager[types.ValidMeasurementUnit](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validMeasurementUnit *types.ValidMeasurementUnit
		validMeasurementUnit, err = dataManager.GetValidMeasurementUnit(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid measurement unit")
		}

		toBeIndexed = converters.ConvertValidMeasurementUnitToValidMeasurementUnitSearchSubset(validMeasurementUnit)
	case IndexTypeValidPreparations:
		im, err = config.ProvideIndexManager[types.ValidPreparation](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validPreparation *types.ValidPreparation
		validPreparation, err = dataManager.GetValidPreparation(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid preparation")
		}

		toBeIndexed = converters.ConvertValidPreparationToValidPreparationSearchSubset(validPreparation)
	case IndexTypeValidIngredientStates:
		im, err = config.ProvideIndexManager[types.ValidIngredientState](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validIngredientState *types.ValidIngredientState
		validIngredientState, err = dataManager.GetValidIngredientState(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient state")
		}

		toBeIndexed = converters.ConvertValidIngredientStateToValidIngredientStateSearchSubset(validIngredientState)
	case IndexTypeValidIngredientMeasurementUnits:
		im, err = config.ProvideIndexManager[types.ValidIngredientMeasurementUnit](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit
		validIngredientMeasurementUnit, err = dataManager.GetValidIngredientMeasurementUnit(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient measurement unit")
		}

		toBeIndexed = converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitSearchSubset(validIngredientMeasurementUnit)
	case IndexTypeValidMeasurementUnitConversions:
		im, err = config.ProvideIndexManager[types.ValidMeasurementUnitConversion](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validMeasurementUnitConversion *types.ValidMeasurementUnitConversion
		validMeasurementUnitConversion, err = dataManager.GetValidMeasurementUnitConversion(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid measurement unit conversion")
		}

		toBeIndexed = converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionSearchSubset(validMeasurementUnitConversion)
	case IndexTypeValidPreparationInstruments:
		im, err = config.ProvideIndexManager[types.ValidPreparationInstrument](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validPreparationInstrument *types.ValidPreparationInstrument
		validPreparationInstrument, err = dataManager.GetValidPreparationInstrument(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid preparation instrument")
		}

		toBeIndexed = converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentSearchSubset(validPreparationInstrument)
	case IndexTypeValidIngredientPreparations:
		im, err = config.ProvideIndexManager[types.ValidIngredientPreparation](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var validIngredientPreparation *types.ValidIngredientPreparation
		validIngredientPreparation, err = dataManager.GetValidIngredientPreparation(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient preparation")
		}

		toBeIndexed = converters.ConvertValidIngredientPreparationToValidIngredientPreparationSearchSubset(validIngredientPreparation)
	default:
		logger.Info("invalid index type specified, exiting")
		return nil
	}

	if searchIndexRequest.Delete {
		if err = im.Delete(ctx, searchIndexRequest.RowID); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "deleting data")
		}

		return nil
	} else {
		if err = im.Index(ctx, searchIndexRequest.RowID, toBeIndexed); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "indexing data")
		}
	}

	return nil
}

type IndexRequest struct {
	RowID     string `json:"rowID"`
	IndexType string `json:"type"`
	Delete    bool   `json:"delete"`
}

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
)

const (
	IndexTypeRecipes                         = "recipes"
	IndexTypeMeals                           = "meals"
	IndexTypeValidIngredients                = "valid_ingredients"
	IndexTypeValidInstruments                = "valid_instruments"
	IndexTypeValidMeasurementUnits           = "valid_measurement_units"
	IndexTypeValidPreparations               = "valid_preparations"
	IndexTypeValidIngredientStates           = "valid_ingredient_states"
	IndexTypeValidIngredientMeasurementUnits = "valid_ingredient_measurement_units"
	IndexTypeValidMeasurementUnitConversions = "valid_measurement_unit_conversions"
	IndexTypeValidPreparationInstruments     = "valid_preparation_instruments"
	IndexTypeValidIngredientPreparations     = "valid_ingredient_preparations"
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
		im, err = config.ProvideIndexManager[search.RecipeSearchSubset](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var recipe *types.Recipe
		recipe, err = dataManager.GetRecipe(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting meal")
		}

		toBeIndexed = SubsetFromRecipe(recipe)
		if err = im.Index(ctx, recipe.ID, toBeIndexed); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "indexing meal")
		}

		return nil
	case IndexTypeMeals:
		im, err = config.ProvideIndexManager[search.MealSearchSubset](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var meal *types.Meal
		meal, err = dataManager.GetMeal(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting meal")
		}

		toBeIndexed = SubsetFromMeal(meal)
	case IndexTypeValidIngredients:
		im, err = config.ProvideIndexManager[types.ValidIngredient](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		toBeIndexed, err = dataManager.GetValidIngredient(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient")
		}
	case IndexTypeValidInstruments:
		im, err = config.ProvideIndexManager[types.ValidInstrument](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		toBeIndexed, err = dataManager.GetValidInstrument(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid instrument")
		}
	case IndexTypeValidMeasurementUnits:
		im, err = config.ProvideIndexManager[types.ValidMeasurementUnit](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		toBeIndexed, err = dataManager.GetValidMeasurementUnit(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid measurement unit")
		}
	case IndexTypeValidPreparations:
		im, err = config.ProvideIndexManager[types.ValidPreparation](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		toBeIndexed, err = dataManager.GetValidPreparation(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid preparation")
		}
	case IndexTypeValidIngredientStates:
		im, err = config.ProvideIndexManager[types.ValidIngredientState](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		toBeIndexed, err = dataManager.GetValidIngredientState(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient state")
		}
	case IndexTypeValidIngredientMeasurementUnits:
		im, err = config.ProvideIndexManager[types.ValidIngredientMeasurementUnit](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		toBeIndexed, err = dataManager.GetValidIngredientMeasurementUnit(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient measurement unit")
		}
	case IndexTypeValidMeasurementUnitConversions:
		im, err = config.ProvideIndexManager[types.ValidMeasurementUnitConversion](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		toBeIndexed, err = dataManager.GetValidMeasurementConversion(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid measurement unit conversion")
		}
	case IndexTypeValidPreparationInstruments:
		im, err = config.ProvideIndexManager[types.ValidPreparationInstrument](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}
		toBeIndexed, err = dataManager.GetValidPreparationInstrument(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid preparation instrument")
		}
	case IndexTypeValidIngredientPreparations:
		im, err = config.ProvideIndexManager[types.ValidIngredientPreparation](ctx, logger, tracerProvider, searchConfig, searchIndexRequest.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		toBeIndexed, err = dataManager.GetValidIngredientPreparation(ctx, searchIndexRequest.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting valid ingredient preparation")
		}
	default:
		logger.Info("invalid index type specified, exiting")
		return nil
	}

	if err = im.Index(ctx, searchIndexRequest.RowID, toBeIndexed); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "indexing meal")
	}

	return nil
}

func SubsetFromRecipe(r *types.Recipe) *search.RecipeSearchSubset {
	x := &search.RecipeSearchSubset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, step := range r.Steps {
		stepSubset := search.RecipeStepSearchSubset{
			Preparation: step.Preparation.Name,
		}

		for _, ingredient := range step.Ingredients {
			stepSubset.Ingredients = append(stepSubset.Ingredients, ingredient.Name)
		}

		for _, instrument := range step.Instruments {
			stepSubset.Instruments = append(stepSubset.Instruments, instrument.Name)
		}

		for _, vessel := range step.Vessels {
			stepSubset.Vessels = append(stepSubset.Vessels, vessel.Name)
		}

		x.Steps = append(x.Steps, stepSubset)
	}

	return x
}

func SubsetFromMeal(r *types.Meal) *search.MealSearchSubset {
	x := &search.MealSearchSubset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, component := range r.Components {
		x.Recipes = append(x.Recipes, component.Recipe.Name)
	}

	return x
}

type IndexRequest struct {
	Value     any    `json:"any"`
	RowID     string `json:"rowID"`
	IndexType string `json:"type"`
}

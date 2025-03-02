package managers

import (
	"context"
	"errors"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
)

const (
	recipeManagerName = "recipe_manager"
)

type (
	RecipeManager interface {
		ListRecipes(ctx context.Context, filter *filtering.QueryFilter) ([]*types.Recipe, string, error)
		CreateRecipe(ctx context.Context, input *types.RecipeCreationRequestInput) (*types.Recipe, error)
		ReadRecipe(ctx context.Context, recipeID string) (*types.Recipe, error)
		SearchRecipes(ctx context.Context, query string, filter *filtering.QueryFilter) ([]*types.Recipe, string, error)
		UpdateRecipe(ctx context.Context, recipeID string, input *types.RecipeUpdateRequestInput) error
		ArchiveRecipe(ctx context.Context, recipeID, ownerID string) error
		RecipeEstimatedPrepSteps(ctx context.Context, recipeID string) ([]*types.MealPlanTaskDatabaseCreationEstimate, error)
		RecipeImageUpload(ctx context.Context)
		RecipeMermaid(ctx context.Context, recipeID string) (string, error)
		CloneRecipe(ctx context.Context, recipeID, newOwnerID string) (*types.Recipe, error)

		ListRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeStep, string, error)
		CreateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepCreationRequestInput) (*types.RecipeStep, error)
		ReadRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error)
		UpdateRecipeStep(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepUpdateRequestInput) error
		ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error
		RecipeStepImageUpload(ctx context.Context)

		ListRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepProduct, string, error)
		CreateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductCreationRequestInput) (*types.RecipeStepProduct, error)
		ReadRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string, input *types.RecipeStepProductUpdateRequestInput) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) error

		ListRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepInstrument, string, error)
		CreateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentCreationRequestInput) (*types.RecipeStepInstrument, error)
		ReadRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error)
		UpdateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string, input *types.RecipeStepInstrumentUpdateRequestInput) error
		ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error

		ListRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepIngredient, string, error)
		CreateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientCreationRequestInput) (*types.RecipeStepIngredient, error)
		ReadRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string, input *types.RecipeStepIngredientUpdateRequestInput) error
		ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error

		ListRecipePrepTask(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipePrepTask, string, error)
		CreateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskCreationRequestInput) (*types.RecipePrepTask, error)
		ReadRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error)
		UpdateRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string, input *types.RecipePrepTaskUpdateRequestInput) error
		ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error

		ListRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepCompletionCondition, string, error)
		CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput) (*types.RecipeStepCompletionCondition, error)
		ReadRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*types.RecipeStepCompletionCondition, error)
		UpdateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string, input *types.RecipeStepCompletionConditionUpdateRequestInput) error
		ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) error

		ListRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepVessel, string, error)
		CreateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselCreationRequestInput) (*types.RecipeStepVessel, error)
		ReadRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error)
		UpdateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string, input *types.RecipeStepVesselUpdateRequestInput) error
		ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) error

		ListRecipeRatings(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeRating, string, error)
		ReadRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*types.RecipeRating, error)
		CreateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingCreationRequestInput) (*types.RecipeRating, error)
		UpdateRecipeRating(ctx context.Context, recipeID, recipeRatingID string, input *types.RecipeRatingUpdateRequestInput) error
		ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error
	}

	recipeManager struct {
		db                   database.DataManager
		tracer               tracing.Tracer
		logger               logging.Logger
		dataChangesPublisher messagequeue.Publisher
		recipeSearchIndex    textsearch.IndexSearcher[eatingindexing.RecipeSearchSubset]
		recipeAnalyzer       recipeanalysis.RecipeAnalyzer
	}
)

var errUnimplemented = errors.New("not implemented") // TODO: DELETE ME

func NewRecipeManager(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	db database.DataManager,
	cfg *msgconfig.QueuesConfig,
	publisherProvider messagequeue.PublisherProvider,
	recipeAnalyzer recipeanalysis.RecipeAnalyzer,
	searchConfig *textsearchcfg.Config,
	metricsProvider metrics.Provider,
) (RecipeManager, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("failed to provide publisher for data changes topic: %w", err)
	}

	validIngredientStatesSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.RecipeSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeRecipes)
	if err != nil {
		return nil, fmt.Errorf("failed to provide search index for %s index", eatingindexing.IndexTypeValidIngredientStates)
	}

	m := &recipeManager{
		db:                   db,
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(recipeManagerName)),
		logger:               logging.EnsureLogger(logger).WithName(recipeManagerName),
		dataChangesPublisher: dataChangesPublisher,
		recipeSearchIndex:    validIngredientStatesSearchIndex,
		recipeAnalyzer:       recipeAnalyzer,
	}

	return m, nil
}

/*

TODO list:

- [x] all returned errors have description strings
- [x] all relevant input params are accounted for in logs
- [x] all relevant input params are accounted for in traces
- [x] all pointer inputs have nil checks
- [x] filters are defaulted
- [ ] no more references to `errUnimplemented`
- [ ] all CUD functions fire a data change event
- [ ] unit tests lmfao

*/

func (m *recipeManager) ListRecipes(ctx context.Context, filter *filtering.QueryFilter) ([]*types.Recipe, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipes(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipes")
	}

	return results.Data, "", nil
}

func (m *recipeManager) CreateRecipe(ctx context.Context, input *types.RecipeCreationRequestInput) (*types.Recipe, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span)

	convertedInput, err := converters.ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput(input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "converting recipe input")
	}

	logger = logger.WithValue(keys.RecipeIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipe(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe")
	}

	return created, nil
}

func (m *recipeManager) ReadRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	x, err := m.db.GetRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
	}

	return x, nil
}

func (m *recipeManager) SearchRecipes(ctx context.Context, query string, filter *filtering.QueryFilter) ([]*types.Recipe, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	recipes, err := m.db.SearchForRecipes(ctx, query, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "searching for recipes")
	}

	return recipes.Data, "", nil
}

func (m *recipeManager) UpdateRecipe(ctx context.Context, recipeID string, input *types.RecipeUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	existingRecipe, err := m.db.GetRecipe(ctx, recipeID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe")
	}

	existingRecipe.Update(input)
	if err = m.db.UpdateRecipe(ctx, existingRecipe); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipe(ctx context.Context, recipeID, ownerID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: recipeID,
		keys.UserIDKey:   ownerID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	if err := m.db.ArchiveRecipe(ctx, recipeID, ownerID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe")
	}

	return errUnimplemented
}

func (m *recipeManager) RecipeEstimatedPrepSteps(ctx context.Context, recipeID string) ([]*types.MealPlanTaskDatabaseCreationEstimate, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	/*

		logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
		tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

		x, err := m.db.GetRecipe(ctx, recipeID)
		if  err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
		}

		stepInputs, err := m.recipeAnalyzer.GenerateMealPlanTasksForRecipe(ctx, "", x)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "generating meal plan tasks")
		}

		responseEvents := []*types.MealPlanTaskDatabaseCreationEstimate{}
		for _, input := range stepInputs {
			responseEvents = append(responseEvents, &types.MealPlanTaskDatabaseCreationEstimate{
				CreationExplanation: input.CreationExplanation,
			})
		}

		return responseEvents, nil

	*/

	return nil, errUnimplemented
}

func (m *recipeManager) RecipeImageUpload(ctx context.Context) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return
}

func (m *recipeManager) RecipeMermaid(ctx context.Context, recipeID string) (string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	/*
		logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
		tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

		recipe, err := m.db.GetRecipe(ctx, recipeID)
		if err != nil {
			return "", observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
		}

		return m.recipeAnalyzer.RenderMermaidDiagramForRecipe(ctx, recipe), nil
	*/

	return "", errUnimplemented
}

func cloneRecipe(x *types.Recipe, userID string) *types.RecipeDatabaseCreationInput {
	ingredientProductIndices := map[string]int{}
	instrumentProductIndices := map[string]int{}
	vesselProductIndices := map[string]int{}
	for _, step := range x.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.RecipeStepProductID != nil {
				ingredientProductIndices[ingredient.ID] = x.FindStepIndexByID(x.FindStepForRecipeStepProductID(*ingredient.RecipeStepProductID).ID)
			}
		}

		for _, instrument := range step.Instruments {
			if instrument.RecipeStepProductID != nil {
				instrumentProductIndices[instrument.ID] = x.FindStepIndexByID(x.FindStepForRecipeStepProductID(*instrument.RecipeStepProductID).ID)
			}
		}

		for _, vessel := range step.Vessels {
			if vessel.RecipeStepProductID != nil {
				vesselProductIndices[vessel.ID] = x.FindStepIndexByID(x.FindStepForRecipeStepProductID(*vessel.RecipeStepProductID).ID)
			}
		}
	}

	// clone recipe.
	cloneInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(x)
	cloneInput.CreatedByUser = userID
	// TODO: cloneInput.ClonedFromRecipeID = &x.ID

	cloneInput.ID = identifiers.New()
	for i := range cloneInput.Steps {
		newRecipeStepID := identifiers.New()
		cloneInput.Steps[i].ID = newRecipeStepID
		for j := range cloneInput.Steps[i].Ingredients {
			if index, ok := ingredientProductIndices[x.Steps[i].Ingredients[j].ID]; ok {
				cloneInput.Steps[i].Ingredients[j].ProductOfRecipeStepIndex = pointer.To(uint64(index))
			}
			cloneInput.Steps[i].Ingredients[j].ID = identifiers.New()
			cloneInput.Steps[i].Ingredients[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].Instruments {
			if index, ok := instrumentProductIndices[x.Steps[i].Instruments[j].ID]; ok {
				cloneInput.Steps[i].Instruments[j].ProductOfRecipeStepIndex = pointer.To(uint64(index))
			}
			cloneInput.Steps[i].Instruments[j].ID = identifiers.New()
			cloneInput.Steps[i].Instruments[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].Vessels {
			if index, ok := vesselProductIndices[x.Steps[i].Vessels[j].ID]; ok {
				cloneInput.Steps[i].Vessels[j].ProductOfRecipeStepIndex = pointer.To(uint64(index))
			}
			cloneInput.Steps[i].Vessels[j].ID = identifiers.New()
			cloneInput.Steps[i].Vessels[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].Products {
			cloneInput.Steps[i].Products[j].ID = identifiers.New()
			cloneInput.Steps[i].Products[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].CompletionConditions {
			newCompletionConditionID := identifiers.New()
			cloneInput.Steps[i].CompletionConditions[j].ID = newCompletionConditionID
			cloneInput.Steps[i].CompletionConditions[j].BelongsToRecipeStep = newRecipeStepID
			for k := range cloneInput.Steps[i].CompletionConditions[j].Ingredients {
				cloneInput.Steps[i].CompletionConditions[j].Ingredients[k].ID = identifiers.New()
				cloneInput.Steps[i].CompletionConditions[j].Ingredients[k].BelongsToRecipeStepCompletionCondition = newCompletionConditionID
			}
		}
	}

	// TODO: handle media here eventually

	for i := range cloneInput.PrepTasks {
		newPrepTaskID := identifiers.New()
		cloneInput.PrepTasks[i].ID = newPrepTaskID
		for j := range cloneInput.PrepTasks[i].TaskSteps {
			cloneInput.PrepTasks[i].TaskSteps[j].ID = identifiers.New()
			cloneInput.PrepTasks[i].TaskSteps[j].BelongsToRecipePrepTask = newPrepTaskID
		}
	}

	return cloneInput
}

func (m *recipeManager) CloneRecipe(ctx context.Context, recipeID, newOwnerID string) (*types.Recipe, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey: recipeID,
		"new_owner":      newOwnerID,
	})
	tracing.AttachToSpan(span, keys.UserIDKey, newOwnerID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	original, err := m.db.GetRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe by id")
	}

	newRecipe, err := m.db.CreateRecipe(ctx, cloneRecipe(original, newOwnerID))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating clone of recipe")
	}

	return newRecipe, nil
}

func (m *recipeManager) ListRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeStep, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipeSteps(ctx, recipeID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipe steps")
	}

	return results.Data, "", nil
}

func (m *recipeManager) CreateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepCreationRequestInput) (*types.RecipeStep, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	convertedInput := converters.ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(input)
	logger = logger.WithValue(keys.RecipeStepIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipeStep(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step")
	}

	return created, nil
}

func (m *recipeManager) ReadRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	x, err := m.db.GetRecipeStep(ctx, recipeID, recipeStepID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step")
	}

	return x, nil
}

func (m *recipeManager) UpdateRecipeStep(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	existingRecipeStep, err := m.db.GetRecipeStep(ctx, recipeID, recipeStepID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe step")
	}

	existingRecipeStep.Update(input)
	if err = m.db.UpdateRecipeStep(ctx, existingRecipeStep); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if err := m.db.ArchiveRecipeStep(ctx, recipeID, recipeStepID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step")
	}

	return nil
}

func (m *recipeManager) RecipeStepImageUpload(ctx context.Context) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return
}

func (m *recipeManager) ListRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepProduct, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipeStepProducts(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipe step products")
	}

	return results.Data, "", nil
}

func (m *recipeManager) CreateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductCreationRequestInput) (*types.RecipeStepProduct, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	convertedInput := converters.ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput(input)
	logger = logger.WithValue(keys.RecipeStepProductIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipeStepProduct(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step product")
	}

	return created, nil
}

func (m *recipeManager) ReadRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:            recipeID,
		keys.RecipeStepIDKey:        recipeStepID,
		keys.RecipeStepProductIDKey: recipeStepProductID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	x, err := m.db.GetRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step product")
	}

	return x, nil
}

func (m *recipeManager) UpdateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string, input *types.RecipeStepProductUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:            recipeID,
		keys.RecipeStepIDKey:        recipeStepID,
		keys.RecipeStepProductIDKey: recipeStepProductID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	existingRecipeStepProduct, err := m.db.GetRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe step product")
	}

	existingRecipeStepProduct.Update(input)
	if err = m.db.UpdateRecipeStepProduct(ctx, existingRecipeStepProduct); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:            recipeID,
		keys.RecipeStepIDKey:        recipeStepID,
		keys.RecipeStepProductIDKey: recipeStepProductID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	// TODO: refactor this to include recipe ID
	if err := m.db.ArchiveRecipeStepProduct(ctx, recipeStepID, recipeStepProductID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step product")
	}

	return nil
}

func (m *recipeManager) ListRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepInstrument, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipeStepInstruments(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipe step instruments")
	}

	return results.Data, "", nil
}

func (m *recipeManager) CreateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentCreationRequestInput) (*types.RecipeStepInstrument, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	convertedInput := converters.ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(input)
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipeStepInstrument(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step instrument")
	}

	return created, nil
}

func (m *recipeManager) ReadRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               recipeID,
		keys.RecipeStepIDKey:           recipeStepID,
		keys.RecipeStepInstrumentIDKey: recipeStepInstrumentID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	x, err := m.db.GetRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step instrument")
	}

	return x, nil
}

func (m *recipeManager) UpdateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string, input *types.RecipeStepInstrumentUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               recipeID,
		keys.RecipeStepIDKey:           recipeStepID,
		keys.RecipeStepInstrumentIDKey: recipeStepInstrumentID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	existingRecipeStepInstrument, err := m.db.GetRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe step instrument")
	}

	existingRecipeStepInstrument.Update(input)
	if err = m.db.UpdateRecipeStepInstrument(ctx, existingRecipeStepInstrument); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step instrument")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               recipeID,
		keys.RecipeStepIDKey:           recipeStepID,
		keys.RecipeStepInstrumentIDKey: recipeStepInstrumentID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	// TODO: refactor this to accept recipe ID
	if err := m.db.ArchiveRecipeStepInstrument(ctx, recipeStepID, recipeStepInstrumentID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step instrument")
	}

	return nil
}

func (m *recipeManager) ListRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepIngredient, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipeStepIngredients(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipe step ingredients")
	}

	return results.Data, "", nil
}

func (m *recipeManager) CreateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientCreationRequestInput) (*types.RecipeStepIngredient, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	convertedInput := converters.ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(input)
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipeStepIngredient(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step ingredient")
	}

	return created, nil
}

func (m *recipeManager) ReadRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               recipeID,
		keys.RecipeStepIDKey:           recipeStepID,
		keys.RecipeStepIngredientIDKey: recipeStepIngredientID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	x, err := m.db.GetRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step ingredient")
	}

	return x, nil
}

func (m *recipeManager) UpdateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string, input *types.RecipeStepIngredientUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               recipeID,
		keys.RecipeStepIDKey:           recipeStepID,
		keys.RecipeStepIngredientIDKey: recipeStepIngredientID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	existingRecipeStepIngredient, err := m.db.GetRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe step ingredient")
	}

	existingRecipeStepIngredient.Update(input)
	if err = m.db.UpdateRecipeStepIngredient(ctx, existingRecipeStepIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step ingredient")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:               recipeID,
		keys.RecipeStepIDKey:           recipeStepID,
		keys.RecipeStepIngredientIDKey: recipeStepIngredientID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	// TODO: refactor this to include recipe ID
	if err := m.db.ArchiveRecipeStepIngredient(ctx, recipeStepID, recipeStepIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step ingredient")
	}

	return nil
}

func (m *recipeManager) ListRecipePrepTask(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipePrepTask, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipePrepTasksForRecipe(ctx, recipeID)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipe prep tasks")
	}

	return results, "", nil
}

func (m *recipeManager) CreateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskCreationRequestInput) (*types.RecipePrepTask, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	convertedInput := converters.ConvertRecipePrepTaskCreationRequestInputToRecipePrepTaskDatabaseCreationInput(input)
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipePrepTask(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
	}

	return created, nil
}

func (m *recipeManager) ReadRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:         recipeID,
		keys.RecipePrepTaskIDKey: recipePrepTaskID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

	x, err := m.db.GetRecipePrepTask(ctx, recipeID, recipePrepTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe prep task")
	}

	return x, nil
}

func (m *recipeManager) UpdateRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string, input *types.RecipePrepTaskUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:         recipeID,
		keys.RecipePrepTaskIDKey: recipePrepTaskID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

	existingRecipePrepTask, err := m.db.GetRecipePrepTask(ctx, recipeID, recipePrepTaskID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe prep task")
	}

	existingRecipePrepTask.Update(input)
	if err = m.db.UpdateRecipePrepTask(ctx, existingRecipePrepTask); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:         recipeID,
		keys.RecipePrepTaskIDKey: recipePrepTaskID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

	if err := m.db.ArchiveRecipePrepTask(ctx, recipeID, recipePrepTaskID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe prep task")
	}

	return nil
}

func (m *recipeManager) ListRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepCompletionCondition, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipeStepCompletionConditions(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipe step completion conditions")
	}

	return results.Data, "", nil
}

func (m *recipeManager) CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput) (*types.RecipeStepCompletionCondition, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	convertedInput := converters.ConvertRecipeStepCompletionConditionForExistingRecipeCreationRequestInputToRecipeStepCompletionConditionDatabaseCreationInput(input)
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipeStepCompletionCondition(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step completion condition")
	}

	return created, nil
}

func (m *recipeManager) ReadRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*types.RecipeStepCompletionCondition, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:                        recipeID,
		keys.RecipeStepIDKey:                    recipeStepID,
		keys.RecipeStepCompletionConditionIDKey: recipeStepCompletionConditionID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	x, err := m.db.GetRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step completion condition")
	}

	return x, nil
}

func (m *recipeManager) UpdateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string, input *types.RecipeStepCompletionConditionUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:                        recipeID,
		keys.RecipeStepIDKey:                    recipeStepID,
		keys.RecipeStepCompletionConditionIDKey: recipeStepCompletionConditionID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	existingRecipeStepCompletionCondition, err := m.db.GetRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe step completion condition")
	}

	existingRecipeStepCompletionCondition.Update(input)
	if err = m.db.UpdateRecipeStepCompletionCondition(ctx, existingRecipeStepCompletionCondition); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step completion condition")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:                        recipeID,
		keys.RecipeStepIDKey:                    recipeStepID,
		keys.RecipeStepCompletionConditionIDKey: recipeStepCompletionConditionID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	// TODO: refactor this to include recipe ID
	if err := m.db.ArchiveRecipeStepCompletionCondition(ctx, recipeStepID, recipeStepCompletionConditionID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step completion condition")
	}

	return nil
}

func (m *recipeManager) ListRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepVessel, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipeStepVessels(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipe step vessels")
	}

	return results.Data, "", nil
}

func (m *recipeManager) CreateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselCreationRequestInput) (*types.RecipeStepVessel, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	convertedInput := converters.ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput(input)
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipeStepVessel(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step vessel")
	}

	return created, nil
}

func (m *recipeManager) ReadRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:           recipeID,
		keys.RecipeStepIDKey:       recipeStepID,
		keys.RecipeStepVesselIDKey: recipeStepVesselID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)

	x, err := m.db.GetRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step vessel")
	}

	return x, nil
}

func (m *recipeManager) UpdateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string, input *types.RecipeStepVesselUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:           recipeID,
		keys.RecipeStepIDKey:       recipeStepID,
		keys.RecipeStepVesselIDKey: recipeStepVesselID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)

	existingRecipeStepVessel, err := m.db.GetRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe step vessel")
	}

	existingRecipeStepVessel.Update(input)
	if err = m.db.UpdateRecipeStepVessel(ctx, existingRecipeStepVessel); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step vessel")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:           recipeID,
		keys.RecipeStepIDKey:       recipeStepID,
		keys.RecipeStepVesselIDKey: recipeStepVesselID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)

	if err := m.db.ArchiveRecipeStepVessel(ctx, recipeStepID, recipeStepVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step vessel")
	}

	return nil
}

func (m *recipeManager) ListRecipeRatings(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeRating, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	results, err := m.db.GetRecipeRatingsForRecipe(ctx, recipeID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching list of recipe ratings")
	}

	return results.Data, "", nil
}

func (m *recipeManager) ReadRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*types.RecipeRating, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:       recipeID,
		keys.RecipeRatingIDKey: recipeRatingID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	x, err := m.db.GetRecipeRating(ctx, recipeID, recipeRatingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe rating")
	}

	return x, nil
}

func (m *recipeManager) CreateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingCreationRequestInput) (*types.RecipeRating, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	convertedInput := converters.ConvertRecipeRatingCreationRequestInputToRecipeRatingDatabaseCreationInput(input)
	logger = logger.WithValue(keys.RecipeRatingIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, convertedInput.ID)

	created, err := m.db.CreateRecipeRating(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe rating")
	}

	return created, nil
}

func (m *recipeManager) UpdateRecipeRating(ctx context.Context, recipeID, recipeRatingID string, input *types.RecipeRatingUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:       recipeID,
		keys.RecipeRatingIDKey: recipeRatingID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	existingRecipeRating, err := m.db.GetRecipeRating(ctx, recipeID, recipeRatingID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving existing recipe rating")
	}

	existingRecipeRating.Update(input)
	if err = m.db.UpdateRecipeRating(ctx, existingRecipeRating); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe rating")
	}

	return nil
}

func (m *recipeManager) ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:       recipeID,
		keys.RecipeRatingIDKey: recipeRatingID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	if err := m.db.ArchiveRecipeRating(ctx, recipeID, recipeRatingID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe rating")
	}

	return nil
}

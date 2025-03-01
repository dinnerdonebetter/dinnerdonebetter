package managers

import (
	"context"
	"errors"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
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
		UpdateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepUpdateRequestInput) error
		ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error
		RecipeStepImageUpload(ctx context.Context)

		ListRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepProduct, string, error)
		CreateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductCreationRequestInput) (*types.RecipeStepProduct, error)
		ReadRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductUpdateRequestInput) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) error

		ListRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepInstrument, string, error)
		CreateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentCreationRequestInput) (*types.RecipeStepInstrument, error)
		ReadRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error)
		UpdateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentUpdateRequestInput) error
		ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error

		ListRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepIngredient, string, error)
		CreateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientCreationRequestInput) (*types.RecipeStepIngredient, error)
		ReadRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientUpdateRequestInput) error
		ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error

		ListRecipePrepTask(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipePrepTask, string, error)
		CreateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskCreationRequestInput) (*types.RecipePrepTask, error)
		ReadRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error)
		UpdateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskUpdateRequestInput) error
		ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error

		ListRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepCompletionCondition, string, error)
		CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionCreationRequestInput) (*types.RecipeStepCompletionCondition, error)
		ReadRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*types.RecipeStepCompletionCondition, error)
		UpdateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionUpdateRequestInput) error
		ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) error

		ListRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepVessel, string, error)
		CreateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselCreationRequestInput) (*types.RecipeStepVessel, error)
		ReadRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error)
		UpdateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselUpdateRequestInput) error
		ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) error

		ListRecipeRatings(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeRating, string, error)
		ReadRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*types.RecipeRating, error)
		CreateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingCreationRequestInput) (*types.RecipeRating, error)
		UpdateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingUpdateRequestInput) error
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
	}

	return m, nil
}

/*

TODO list:

- [ ] all returned errors have description strings
- [ ] all relevant input params are accounted for in logs
- [ ] all relevant input params are accounted for in traces
- [ ] all pointer inputs have nil checks
- [ ] filters are defaulted
- [ ] no more references to `errUnimplemented`
- [ ] all query filters are defaulted when nil
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

	recipe, err := m.db.GetRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
	}

	return recipe, nil
}

func (m *recipeManager) SearchRecipes(ctx context.Context, query string, filter *filtering.QueryFilter) ([]*types.Recipe, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	recipes, err := m.db.SearchForRecipes(ctx, query, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "searching for recipes")
	}

	return recipes.Data, "", nil
}

func (m *recipeManager) UpdateRecipe(ctx context.Context, recipeID string, input *types.RecipeUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

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

	logger := m.logger.WithSpan(span)

	recipeSteps, err := m.db.GetRecipeSteps(ctx, recipeID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing recipe steps")
	}

	return recipeSteps.Data, "", nil
}

func (m *recipeManager) CreateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepCreationRequestInput) (*types.RecipeStep, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.CreateRecipeStep(ctx)

	return nil, errUnimplemented
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

	m.db.GetRecipeStep(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) UpdateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.UpdateRecipeStep(ctx)

	return errUnimplemented
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

	m.db.ArchiveRecipeStep(ctx)

	return errUnimplemented
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

	m.db.GetRecipeStepProducts(ctx)

	return nil, "", errUnimplemented
}

func (m *recipeManager) CreateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductCreationRequestInput) (*types.RecipeStepProduct, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.CreateRecipeStepProduct(ctx)

	return nil, errUnimplemented
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

	m.db.GetRecipeStepProduct(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) UpdateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.UpdateRecipeStepProduct(ctx)

	return errUnimplemented
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

	m.db.ArchiveRecipeStepProduct(ctx)

	return errUnimplemented
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

	m.db.GetRecipeStepInstruments(ctx)

	return nil, "", errUnimplemented
}

func (m *recipeManager) CreateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentCreationRequestInput) (*types.RecipeStepInstrument, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.CreateRecipeStepInstrument(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) ReadRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.GetRecipeStepInstrument(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) UpdateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.UpdateRecipeStepInstrument(ctx)

	return errUnimplemented
}

func (m *recipeManager) ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.ArchiveRecipeStepInstrument(ctx)

	return errUnimplemented
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

	m.db.GetRecipeStepIngredients(ctx)

	return nil, "", errUnimplemented
}

func (m *recipeManager) CreateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientCreationRequestInput) (*types.RecipeStepIngredient, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.CreateRecipeStepIngredient(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) ReadRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.GetRecipeStepIngredient(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) UpdateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.UpdateRecipeStepIngredient(ctx)

	return errUnimplemented
}

func (m *recipeManager) ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.ArchiveRecipeStepIngredient(ctx)

	return errUnimplemented
}

func (m *recipeManager) ListRecipePrepTask(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipePrepTask, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.GetRecipePrepTask(ctx)

	return nil, "", errUnimplemented
}

func (m *recipeManager) CreateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskCreationRequestInput) (*types.RecipePrepTask, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.CreateRecipePrepTask(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) ReadRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.GetRecipePrepTask(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) UpdateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.UpdateRecipePrepTask(ctx)

	return errUnimplemented
}

func (m *recipeManager) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.ArchiveRecipePrepTask(ctx)

	return errUnimplemented
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

	m.db.GetRecipeStepCompletionConditions(ctx)

	return nil, "", errUnimplemented
}

func (m *recipeManager) CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionCreationRequestInput) (*types.RecipeStepCompletionCondition, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.CreateRecipeStepCompletionCondition(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) ReadRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*types.RecipeStepCompletionCondition, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.GetRecipeStepCompletionCondition(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) UpdateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.UpdateRecipeStepCompletionCondition(ctx)

	return errUnimplemented
}

func (m *recipeManager) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.ArchiveRecipeStepCompletionCondition(ctx)

	return errUnimplemented
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

	m.db.GetRecipeStepVessels(ctx)

	return nil, "", errUnimplemented
}

func (m *recipeManager) CreateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselCreationRequestInput) (*types.RecipeStepVessel, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.CreateRecipeStepVessel(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) ReadRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.GetRecipeStepVessel(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) UpdateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.RecipeIDKey:     recipeID,
		keys.RecipeStepIDKey: recipeStepID,
	})
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	m.db.UpdateRecipeStepVessel(ctx)

	return errUnimplemented
}

func (m *recipeManager) ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.ArchiveRecipeStepVessel(ctx)

	return errUnimplemented
}

func (m *recipeManager) ListRecipeRatings(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeRating, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.GetRecipeRatings(ctx)

	return nil, "", errUnimplemented
}

func (m *recipeManager) ReadRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*types.RecipeRating, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.GetRecipeRating(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) CreateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingCreationRequestInput) (*types.RecipeRating, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.CreateRecipeRating(ctx)

	return nil, errUnimplemented
}

func (m *recipeManager) UpdateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.UpdateRecipeRating(ctx)

	return errUnimplemented
}

func (m *recipeManager) ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	m.db.ArchiveRecipeRating(ctx)

	return errUnimplemented
}

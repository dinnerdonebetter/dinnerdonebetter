package workers

import (
	"context"

	"github.com/hashicorp/go-multierror"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/recipeanalysis"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	advancedPrepStepCreationEnsurerWorkerName = "advanced_prep_step_creation_ensurer"
)

// AdvancedPrepStepCreationEnsurerWorker ensurers advanced prep steps are created.
type AdvancedPrepStepCreationEnsurerWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	analyzer              recipeanalysis.RecipeAnalyzer
	encoder               encoding.ClientEncoder
	dataManager           database.DataManager
	postUpdatesPublisher  messagequeue.Publisher
	customerDataCollector customerdata.Collector
}

// ProvideAdvancedPrepStepCreationEnsurerWorker provides a AdvancedPrepStepCreationEnsurerWorker.
func ProvideAdvancedPrepStepCreationEnsurerWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	grapher recipeanalysis.RecipeAnalyzer,
	postUpdatesPublisher messagequeue.Publisher,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *AdvancedPrepStepCreationEnsurerWorker {
	return &AdvancedPrepStepCreationEnsurerWorker{
		logger:                logging.EnsureLogger(logger).WithName(advancedPrepStepCreationEnsurerWorkerName),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(advancedPrepStepCreationEnsurerWorkerName)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		analyzer:              grapher,
		postUpdatesPublisher:  postUpdatesPublisher,
		customerDataCollector: customerDataCollector,
	}
}

// HandleMessage handles a pending write.
func (w *AdvancedPrepStepCreationEnsurerWorker) HandleMessage(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	optionsAndSteps, err := w.DetermineCreatableSteps(ctx)
	if err != nil {
		return observability.PrepareError(err, nil, "determining creatable steps")
	}

	logger = logger.WithValue("creatable_steps_qty", len(optionsAndSteps))

	var result *multierror.Error
	for mealPlanOptionID, steps := range optionsAndSteps {
		l := logger.Clone().WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID).WithValue("creatable_prep_step_qty", len(steps))

		createdSteps, creationErr := w.dataManager.CreateAdvancedPrepStepsForMealPlanOption(ctx, mealPlanOptionID, steps)
		if creationErr != nil {
			result = multierror.Append(result, creationErr)
			observability.AcknowledgeError(creationErr, l, span, "creating advanced prep steps for meal plan optino")
		}

		for _, createdStep := range createdSteps {
			if publishErr := w.postUpdatesPublisher.Publish(ctx, &types.DataChangeMessage{
				DataType:                  types.AdvancedPrepStepDataType,
				EventType:                 types.AdvancedPrepStepCreatedCustomerEventType,
				AdvancedPrepStep:          createdStep,
				AdvancedPrepStepID:        createdStep.ID,
				HouseholdID:               "",
				Context:                   nil,
				AttributableToHouseholdID: "",
			}); publishErr != nil {
				observability.AcknowledgeError(publishErr, l, span, "publishing data change event")
			}
		}
	}

	if result == nil {
		return nil
	}

	return result
}

// DetermineCreatableSteps determines which advanced prep steps are creatable for a recipe.
func (w *AdvancedPrepStepCreationEnsurerWorker) DetermineCreatableSteps(ctx context.Context) (map[string][]*types.AdvancedPrepStepDatabaseCreationInput, error) {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()
	logger.Info("fetching finalized meal plan IDs to determine creatable steps")

	results, err := w.dataManager.GetFinalizedMealPlanIDsForTheNextWeek(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting finalized meal plan data for the next week")
	}

	logger = logger.WithValue("steps_to_create", len(results))
	logger.Info("determining creatable steps")

	inputs := map[string][]*types.AdvancedPrepStepDatabaseCreationInput{}
	for _, result := range results {
		l := logger.Clone().WithValues(map[string]interface{}{
			keys.MealPlanIDKey:       result.MealPlanID,
			keys.MealPlanEventIDKey:  result.MealPlanEventID,
			keys.MealPlanOptionIDKey: result.MealPlanOptionID,
			keys.MealIDKey:           result.MealID,
			"recipe_ids":             result.RecipeIDs,
		})
		l.Info("fetching meal plan event")

		mealPlanEvent, getMealPlanEventErr := w.dataManager.GetMealPlanEvent(ctx, result.MealPlanID, result.MealPlanEventID)
		if getMealPlanEventErr != nil {
			return nil, observability.PrepareAndLogError(getMealPlanEventErr, l, span, "fetching meal plan event")
		}

		if _, ok := inputs[result.MealPlanOptionID]; !ok {
			inputs[result.MealPlanOptionID] = []*types.AdvancedPrepStepDatabaseCreationInput{}
		}

		for _, recipeID := range result.RecipeIDs {
			recipe, getRecipeErr := w.dataManager.GetRecipe(ctx, recipeID)
			if getRecipeErr != nil {
				return nil, observability.PrepareAndLogError(getRecipeErr, l, span, "fetching recipe")
			}

			creatableSteps, determineStepsErr := w.analyzer.GenerateAdvancedStepCreationForRecipe(ctx, mealPlanEvent, result.MealPlanOptionID, recipe)
			if determineStepsErr != nil {
				return nil, observability.PrepareAndLogError(determineStepsErr, l, span, "fetching recipe")
			}

			inputs[result.MealPlanOptionID] = append(inputs[result.MealPlanOptionID], creatableSteps...)
		}
	}

	return inputs, nil
}

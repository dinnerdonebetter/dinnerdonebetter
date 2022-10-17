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
	mealPlanTaskCreationEnsurerWorkerName = "meal_plan_task_creation_ensurer"
)

// MealPlanTaskCreatorWorker ensurers meal plan tasks are created.
type MealPlanTaskCreatorWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	analyzer              recipeanalysis.RecipeAnalyzer
	encoder               encoding.ClientEncoder
	dataManager           database.DataManager
	postUpdatesPublisher  messagequeue.Publisher
	customerDataCollector customerdata.Collector
}

// ProvideMealPlanTaskCreationEnsurerWorker provides a MealPlanTaskCreatorWorker.
func ProvideMealPlanTaskCreationEnsurerWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	grapher recipeanalysis.RecipeAnalyzer,
	postUpdatesPublisher messagequeue.Publisher,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *MealPlanTaskCreatorWorker {
	return &MealPlanTaskCreatorWorker{
		logger:                logging.EnsureLogger(logger).WithName(mealPlanTaskCreationEnsurerWorkerName),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(mealPlanTaskCreationEnsurerWorkerName)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		analyzer:              grapher,
		postUpdatesPublisher:  postUpdatesPublisher,
		customerDataCollector: customerDataCollector,
	}
}

// HandleMessage handles a pending write.
func (w *MealPlanTaskCreatorWorker) HandleMessage(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	mealPlansAndSteps, err := w.DetermineCreatableMealPlanTasks(ctx)
	if err != nil {
		return observability.PrepareError(err, nil, "determining creatable steps")
	}

	logger = logger.WithValue("creatable_steps_qty", len(mealPlansAndSteps))

	var result *multierror.Error
	for mealPlanID, steps := range mealPlansAndSteps {
		l := logger.Clone().WithValue(keys.MealPlanIDKey, mealPlanID).WithValue("creatable_prep_step_qty", len(steps))

		createdSteps, creationErr := w.dataManager.CreateMealPlanTasksForMealPlanOption(ctx, steps)
		if creationErr != nil {
			result = multierror.Append(result, creationErr)
			observability.AcknowledgeError(creationErr, l, span, "creating meal plan tasks for meal plan option")
		}

		for _, createdStep := range createdSteps {
			if publishErr := w.postUpdatesPublisher.Publish(ctx, &types.DataChangeMessage{
				DataType:       types.MealPlanTaskDataType,
				EventType:      types.MealPlanTaskCreatedCustomerEventType,
				MealPlanTask:   createdStep,
				MealPlanTaskID: createdStep.ID,
				Context:        nil,
				// TODO: attribute these
				HouseholdID:               "",
				AttributableToHouseholdID: "",
			}); publishErr != nil {
				observability.AcknowledgeError(publishErr, l, span, "publishing data change event")
			}
		}

		if err = w.dataManager.MarkMealPlanAsHavingTasksCreated(ctx, mealPlanID); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if result == nil {
		return nil
	}

	return result
}

// DetermineCreatableMealPlanTasks determines which meal plan tasks are creatable for a recipe.
func (w *MealPlanTaskCreatorWorker) DetermineCreatableMealPlanTasks(ctx context.Context) (map[string][]*types.MealPlanTaskDatabaseCreationInput, error) {
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

	inputs := map[string][]*types.MealPlanTaskDatabaseCreationInput{}
	for _, result := range results {
		l := logger.Clone().WithValues(map[string]interface{}{
			keys.MealPlanIDKey:       result.MealPlanID,
			keys.MealPlanEventIDKey:  result.MealPlanEventID,
			keys.MealPlanOptionIDKey: result.MealPlanOptionID,
			keys.MealIDKey:           result.MealID,
			"recipe_ids":             result.RecipeIDs,
		})
		l.Info("fetching meal plan event")

		if _, ok := inputs[result.MealPlanID]; !ok {
			inputs[result.MealPlanID] = []*types.MealPlanTaskDatabaseCreationInput{}
		}

		for _, recipeID := range result.RecipeIDs {
			recipe, getRecipeErr := w.dataManager.GetRecipe(ctx, recipeID)
			if getRecipeErr != nil {
				return nil, observability.PrepareAndLogError(getRecipeErr, l, span, "fetching recipe")
			}

			creatableSteps, determineStepsErr := w.analyzer.GenerateMealPlanTasksForRecipe(ctx, result.MealPlanOptionID, recipe)
			if determineStepsErr != nil {
				return nil, observability.PrepareAndLogError(determineStepsErr, l, span, "fetching recipe")
			}

			inputs[result.MealPlanID] = append(inputs[result.MealPlanID], creatableSteps...)
		}
	}

	return inputs, nil
}

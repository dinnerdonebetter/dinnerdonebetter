package workers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/graphing"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	advancedPrepStepCreationEnsurerWorkerName = "advanced_prep_step_creation_ensurer"
)

// AdvancedPrepStepCreationEnsurerWorker ensurers advanced prep steps are created.
type AdvancedPrepStepCreationEnsurerWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	grapher               graphing.RecipeGrapher
	encoder               encoding.ClientEncoder
	dataManager           database.DataManager
	postUpdatesPublisher  messagequeue.Publisher
	customerDataCollector customerdata.Collector
}

// ProvideAdvancedPrepStepCreationEnsurerWorker provides a AdvancedPrepStepCreationEnsurerWorker.
func ProvideAdvancedPrepStepCreationEnsurerWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	grapher graphing.RecipeGrapher,
	postUpdatesPublisher messagequeue.Publisher,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *AdvancedPrepStepCreationEnsurerWorker {
	return &AdvancedPrepStepCreationEnsurerWorker{
		logger:                logging.EnsureLogger(logger).WithName(advancedPrepStepCreationEnsurerWorkerName),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(advancedPrepStepCreationEnsurerWorkerName)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		grapher:               grapher,
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
			}); err != nil {
				observability.AcknowledgeError(publishErr, l, span, "publishing data change event")
			}
		}
	}

	if result == nil {
		return nil
	}

	return result
}

// frozenIngredientDefrostStepsFilter iterates through a recipe and returns
// the list of ingredients within that are indicated as kept frozen.
func frozenIngredientDefrostStepsFilter(recipe *types.Recipe) map[string][]int {
	out := map[string][]int{}

	for _, recipeStep := range recipe.Steps {
		ingredientIndices := []int{}
		for i, ingredient := range recipeStep.Ingredients {
			// if it's a valid ingredient
			if ingredient.Ingredient != nil &&
				// if the ingredient has storage temperature set
				ingredient.Ingredient.MinimumIdealStorageTemperatureInCelsius != nil &&
				// the ingredient's storage temperature is set to something about freezing temperature.
				*ingredient.Ingredient.MinimumIdealStorageTemperatureInCelsius <= 3 {
				ingredientIndices = append(ingredientIndices, i)
			}
		}

		if len(ingredientIndices) > 0 {
			out[recipeStep.ID] = ingredientIndices
		}
	}

	return out
}

func whicheverIsLater(t1, t2 time.Time) time.Time {
	if t2.After(t1) {
		return t2
	}
	return t1
}

func buildThawStepCreationExplanation(ingredientIndices []int) string {
	if len(ingredientIndices) == 0 {
		return ""
	}

	stringIndices := []string{}
	for _, i := range ingredientIndices {
		stringIndices = append(stringIndices, fmt.Sprintf("#%d", i+1))
	}

	var d string
	if len(ingredientIndices) > 1 {
		d = "ingredients"
	} else if len(ingredientIndices) == 1 {
		d = "ingredient"
	}

	return fmt.Sprintf("frozen %s (%s) might need to be thawed ahead of time", d, strings.Join(stringIndices, ", "))
}

func determineCreationMinAndMaxTimesForRecipeStep(step *types.RecipeStep, mealPlanEvent *types.MealPlanEvent) (cannotCompleteBefore, cannotCompleteAfter time.Time) {
	var (
		shortestDuration,
		longestDuration uint32
	)

	for _, product := range step.Products {
		if product.MaximumStorageDurationInSeconds < shortestDuration || shortestDuration == 0 {
			shortestDuration = product.MaximumStorageDurationInSeconds
		}

		if product.MaximumStorageDurationInSeconds > longestDuration || longestDuration == 0 {
			longestDuration = product.MaximumStorageDurationInSeconds
		}
	}

	cannotCompleteBefore = whicheverIsLater(time.Now(), mealPlanEvent.StartsAt.Add(time.Duration(shortestDuration)*-time.Second))
	cannotCompleteAfter = whicheverIsLater(mealPlanEvent.StartsAt, mealPlanEvent.StartsAt.Add(time.Duration(longestDuration)*-time.Second))

	return cannotCompleteBefore, cannotCompleteAfter
}

const advancedStepCreationExplanation = "adequate storage instructions for early step"

// DetermineCreatableSteps determines which advanced prep steps are creatable for a recipe.
func (w *AdvancedPrepStepCreationEnsurerWorker) DetermineCreatableSteps(ctx context.Context) (map[string][]*types.AdvancedPrepStepDatabaseCreationInput, error) {
	logger := w.logger.Clone()
	logger.Info("fetching finalized meal plan IDs to determine creatable steps")

	results, err := w.dataManager.GetFinalizedMealPlanIDsForTheNextWeek(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting finalized meal plan data for the next week")
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
			return nil, observability.PrepareAndLogError(getMealPlanEventErr, l, nil, "fetching meal plan event")
		}

		if _, ok := inputs[result.MealPlanOptionID]; !ok {
			inputs[result.MealPlanOptionID] = []*types.AdvancedPrepStepDatabaseCreationInput{}
		}

		for _, recipeID := range result.RecipeIDs {
			recipe, getRecipeErr := w.dataManager.GetRecipe(ctx, recipeID)
			if getRecipeErr != nil {
				return nil, observability.PrepareAndLogError(getRecipeErr, l, nil, "fetching recipe")
			}

			frozenIngredientSteps := frozenIngredientDefrostStepsFilter(recipe)
			logger.WithValue("frozen_steps_qty", len(frozenIngredientSteps)).Info("creating frozen step inputs")

			for stepID, ingredientIndices := range frozenIngredientSteps {
				explanation := buildThawStepCreationExplanation(ingredientIndices)
				if explanation == "" {
					continue
				}

				inputs[result.MealPlanOptionID] = append(inputs[result.MealPlanOptionID], &types.AdvancedPrepStepDatabaseCreationInput{
					ID:                   ksuid.New().String(),
					CannotCompleteBefore: mealPlanEvent.StartsAt.Add(2 * -time.Hour * 24),
					CannotCompleteAfter:  mealPlanEvent.StartsAt.Add(-time.Hour * 24),
					Status:               types.AdvancedPrepStepStatusUnfinished,
					CreationExplanation:  explanation,
					MealPlanOptionID:     result.MealPlanOptionID,
					RecipeStepID:         stepID,
				})
			}

			steps, graphErr := w.grapher.FindStepsEligibleForAdvancedCreation(ctx, recipe)
			if graphErr != nil {
				return nil, observability.PrepareAndLogError(graphErr, l, nil, "generating graph for recipe")
			}

			logger.WithValue("advanced_steps_qty", len(steps)).Info("creating advanced prep step inputs")

			for _, step := range steps {
				cannotCompleteBefore, cannotCompleteAfter := determineCreationMinAndMaxTimesForRecipeStep(step, mealPlanEvent)

				inputs[result.MealPlanOptionID] = append(inputs[result.MealPlanOptionID], &types.AdvancedPrepStepDatabaseCreationInput{
					ID:                   ksuid.New().String(),
					CannotCompleteBefore: cannotCompleteBefore,
					CannotCompleteAfter:  cannotCompleteAfter,
					Status:               types.AdvancedPrepStepStatusUnfinished,
					CreationExplanation:  advancedStepCreationExplanation,
					MealPlanOptionID:     result.MealPlanOptionID,
					RecipeStepID:         step.ID,
				})
			}
		}
	}

	return inputs, nil
}

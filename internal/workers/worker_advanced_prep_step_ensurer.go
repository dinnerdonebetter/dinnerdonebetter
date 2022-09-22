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

	steps, err := w.DetermineCreatableSteps(ctx)
	if err != nil {
		return observability.PrepareError(err, nil, "determining creatable steps")
	}

	var result *multierror.Error
	for _, step := range steps {
		createdStep, creationErr := w.dataManager.CreateAdvancedPrepStep(ctx, step)
		if creationErr != nil {
			result = multierror.Append(result, creationErr)
		}

		if publishErr := w.postUpdatesPublisher.Publish(ctx, &types.DataChangeMessage{
			DataType:                  types.AdvancedPrepStepDataType,
			EventType:                 types.AdvancedPrepStepCreatedCustomerEventType,
			AdvancedPrepStep:          createdStep,
			AdvancedPrepStepID:        createdStep.ID,
			HouseholdID:               "",
			Context:                   nil,
			AttributableToHouseholdID: "",
		}); err != nil {
			result = multierror.Append(result, publishErr)
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

// DetermineCreatableSteps determines which advanced prep steps are creatable for a recipe.
func (w *AdvancedPrepStepCreationEnsurerWorker) DetermineCreatableSteps(ctx context.Context) ([]*types.AdvancedPrepStepDatabaseCreationInput, error) {
	results, err := w.dataManager.GetFinalizedMealPlanIDsForTheNextWeek(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting finalized meal plan data for the next week")
	}

	inputs := []*types.AdvancedPrepStepDatabaseCreationInput{}
	for _, result := range results {
		mealPlanEvent, getMealPlanEventErr := w.dataManager.GetMealPlanEvent(ctx, result.MealPlanID, result.MealPlanEventID)
		if getMealPlanEventErr != nil {
			return nil, observability.PrepareError(getMealPlanEventErr, nil, "fetching meal plan event")
		}

		for _, recipeID := range result.RecipeIDs {
			recipe, getRecipeErr := w.dataManager.GetRecipe(ctx, recipeID)
			if getRecipeErr != nil {
				return nil, observability.PrepareError(getRecipeErr, nil, "fetching recipe")
			}

			steps, graphErr := w.grapher.FindStepsEligibleForAdvancedCreation(ctx, recipe)
			if graphErr != nil {
				return nil, observability.PrepareError(graphErr, nil, "generating graph for recipe")
			}

			frozenIngredientSteps := frozenIngredientDefrostStepsFilter(recipe)
			for stepID, ingredientIndices := range frozenIngredientSteps {
				explanation := buildThawStepCreationExplanation(ingredientIndices)
				if explanation == "" {
					continue
				}

				inputs = append(inputs, &types.AdvancedPrepStepDatabaseCreationInput{
					ID:                   ksuid.New().String(),
					CannotCompleteBefore: mealPlanEvent.StartsAt.Add(2 * -time.Hour * 24),
					CannotCompleteAfter:  mealPlanEvent.StartsAt.Add(-time.Hour * 24),
					Status:               types.AdvancedPrepStepStatusUnfinished,
					CreationExplanation:  explanation,
					MealPlanOptionID:     result.MealPlanOptionID,
					RecipeStepID:         stepID,
				})
			}

			for _, step := range steps {
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

				cannotCompleteBefore := whicheverIsLater(time.Now(), mealPlanEvent.StartsAt.Add(time.Duration(shortestDuration)*-time.Second))
				cannotCompleteAfter := whicheverIsLater(mealPlanEvent.StartsAt, mealPlanEvent.StartsAt.Add(time.Duration(longestDuration)*-time.Second))

				inputs = append(inputs, &types.AdvancedPrepStepDatabaseCreationInput{
					ID:                   ksuid.New().String(),
					CannotCompleteBefore: cannotCompleteBefore,
					CannotCompleteAfter:  cannotCompleteAfter,
					Status:               types.AdvancedPrepStepStatusUnfinished,
					CreationExplanation:  "adequate storage instructions for early step",
					MealPlanOptionID:     result.MealPlanOptionID,
					RecipeStepID:         step.ID,
				})
			}
		}
	}

	return inputs, nil
}

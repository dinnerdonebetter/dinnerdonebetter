package workers

import (
	"context"
	"fmt"
	"time"

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

	_, err := w.DetermineCreatableSteps(ctx)
	return err
}

// frozenIngredientDefrostStepsFilter iterates through a recipe and returns the list of
// ingredients within that are indicated as kept frozen.
func frozenIngredientDefrostStepsFilter(recipe *types.Recipe) []*types.RecipeStep {
	out := []*types.RecipeStep{}

	for _, recipeStep := range recipe.Steps {
		for _, ingredient := range recipeStep.Ingredients {
			// if it's a valid ingredient
			if ingredient.Ingredient != nil &&
				// if the ingredient has storage temperature set
				ingredient.Ingredient.MinimumIdealStorageTemperatureInCelsius != nil &&
				// the ingredient's storage temperature is set to something about freezing temperature.
				*ingredient.Ingredient.MinimumIdealStorageTemperatureInCelsius <= 3 {
				out = append(out, recipeStep)
			}
		}
	}

	return out
}

func earlierTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t2
	} else {
		return t1
	}
}

func whicheverIsLater(t1, t2 time.Time) time.Time {
	if t2.After(t1) {
		return t2
	} else {
		return t1
	}
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
			for _, step := range frozenIngredientSteps {
				inputs = append(inputs, &types.AdvancedPrepStepDatabaseCreationInput{
					ID:                   ksuid.New().String(),
					CannotCompleteBefore: mealPlanEvent.StartsAt.Add(2 * -time.Hour * 24),
					CannotCompleteAfter:  mealPlanEvent.StartsAt.Add(-time.Hour * 24),
					Status:               types.AdvancedPrepStepStatusUnfinished,
					CreationExplanation:  "frozen ingredient in need of thawing",
					MealPlanOptionID:     result.MealPlanOptionID,
					RecipeStepID:         step.ID,
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
					CreationExplanation:  "",
					MealPlanOptionID:     result.MealPlanOptionID,
					RecipeStepID:         step.ID,
				})
			}
		}
	}

	return inputs, nil
}

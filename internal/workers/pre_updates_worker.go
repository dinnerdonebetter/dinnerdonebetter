package workers

import (
	"context"
	"fmt"
	"net/http"

	database "github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/encoding"
	publishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	observability "github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/search"
	"github.com/prixfixeco/api_server/pkg/types"
)

// PreUpdatesWorker updates data from the pending updates topic to the database.
type PreUpdatesWorker struct {
	logger                                  logging.Logger
	tracer                                  tracing.Tracer
	encoder                                 encoding.ClientEncoder
	postUpdatesPublisher                    publishers.Publisher
	dataManager                             database.DataManager
	validInstrumentsIndexManager            search.IndexManager
	validIngredientsIndexManager            search.IndexManager
	validPreparationsIndexManager           search.IndexManager
	validIngredientPreparationsIndexManager search.IndexManager
	recipesIndexManager                     search.IndexManager
	recipeStepsIndexManager                 search.IndexManager
	recipeStepInstrumentsIndexManager       search.IndexManager
	recipeStepIngredientsIndexManager       search.IndexManager
	recipeStepProductsIndexManager          search.IndexManager
	mealPlansIndexManager                   search.IndexManager
	mealPlanOptionsIndexManager             search.IndexManager
	mealPlanOptionVotesIndexManager         search.IndexManager
}

// ProvidePreUpdatesWorker provides a PreUpdatesWorker.
func ProvidePreUpdatesWorker(
	ctx context.Context,
	logger logging.Logger,
	client *http.Client,
	dataManager database.DataManager,
	postUpdatesPublisher publishers.Publisher,
	searchIndexLocation search.IndexPath,
	searchIndexProvider search.IndexManagerProvider,
) (*PreUpdatesWorker, error) {
	const name = "pre_updates"

	validInstrumentsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "valid_instruments", "name", "variant", "description", "icon")
	if err != nil {
		return nil, fmt.Errorf("setting up valid instruments search index manager: %w", err)
	}

	validIngredientsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "valid_ingredients", "name", "variant", "description", "warning", "icon")
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredients search index manager: %w", err)
	}

	validPreparationsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "valid_preparations", "name", "description", "icon")
	if err != nil {
		return nil, fmt.Errorf("setting up valid preparations search index manager: %w", err)
	}

	validIngredientPreparationsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "valid_ingredient_preparations", "notes", "validPreparationID", "validIngredientID")
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient preparations search index manager: %w", err)
	}

	recipesIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "recipes", "name", "source", "description", "inspiredByRecipeID")
	if err != nil {
		return nil, fmt.Errorf("setting up recipes search index manager: %w", err)
	}

	recipeStepsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "recipe_steps", "preparationID", "notes", "recipeID")
	if err != nil {
		return nil, fmt.Errorf("setting up recipe steps search index manager: %w", err)
	}

	recipeStepInstrumentsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "recipe_step_instruments", "instrumentID", "recipeStepID", "notes")
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step instruments search index manager: %w", err)
	}

	recipeStepIngredientsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "recipe_step_ingredients", "ingredientID", "quantityType", "quantityNotes", "ingredientNotes")
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step ingredients search index manager: %w", err)
	}

	recipeStepProductsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "recipe_step_products", "name", "recipeStepID")
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step products search index manager: %w", err)
	}

	mealPlansIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "meal_plans", "state")
	if err != nil {
		return nil, fmt.Errorf("setting up meal plans search index manager: %w", err)
	}

	mealPlanOptionsIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "meal_plan_options", "mealPlanID", "recipeID", "notes")
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan options search index manager: %w", err)
	}

	mealPlanOptionVotesIndexManager, err := searchIndexProvider(ctx, logger, client, searchIndexLocation, "meal_plan_option_votes", "mealPlanOptionID", "notes")
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan option votes search index manager: %w", err)
	}

	w := &PreUpdatesWorker{
		logger:                                  logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                                  tracing.NewTracer(name),
		encoder:                                 encoding.ProvideClientEncoder(logger, encoding.ContentTypeJSON),
		postUpdatesPublisher:                    postUpdatesPublisher,
		dataManager:                             dataManager,
		validInstrumentsIndexManager:            validInstrumentsIndexManager,
		validIngredientsIndexManager:            validIngredientsIndexManager,
		validPreparationsIndexManager:           validPreparationsIndexManager,
		validIngredientPreparationsIndexManager: validIngredientPreparationsIndexManager,
		recipesIndexManager:                     recipesIndexManager,
		recipeStepsIndexManager:                 recipeStepsIndexManager,
		recipeStepInstrumentsIndexManager:       recipeStepInstrumentsIndexManager,
		recipeStepIngredientsIndexManager:       recipeStepIngredientsIndexManager,
		recipeStepProductsIndexManager:          recipeStepProductsIndexManager,
		mealPlansIndexManager:                   mealPlansIndexManager,
		mealPlanOptionsIndexManager:             mealPlanOptionsIndexManager,
		mealPlanOptionVotesIndexManager:         mealPlanOptionVotesIndexManager,
	}

	return w, nil
}

// HandleMessage handles a pending update.
func (w *PreUpdatesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var msg *types.PreUpdateMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	logger := w.logger.WithValue("data_type", msg.DataType)

	logger.Debug("message read")

	switch msg.DataType {
	case types.ValidInstrumentDataType:
		if err := w.dataManager.UpdateValidInstrument(ctx, msg.ValidInstrument); err != nil {
			return observability.PrepareError(err, logger, span, "creating valid instrument")
		}

		if err := w.validInstrumentsIndexManager.Index(ctx, msg.ValidInstrument.ID, msg.ValidInstrument); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the valid instrument")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "validInstrumentUpdated",
				ValidInstrument:           msg.ValidInstrument,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.ValidIngredientDataType:
		if err := w.dataManager.UpdateValidIngredient(ctx, msg.ValidIngredient); err != nil {
			return observability.PrepareError(err, logger, span, "creating valid ingredient")
		}

		if err := w.validIngredientsIndexManager.Index(ctx, msg.ValidIngredient.ID, msg.ValidIngredient); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the valid ingredient")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "validIngredientUpdated",
				ValidIngredient:           msg.ValidIngredient,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.ValidPreparationDataType:
		if err := w.dataManager.UpdateValidPreparation(ctx, msg.ValidPreparation); err != nil {
			return observability.PrepareError(err, logger, span, "creating valid preparation")
		}

		if err := w.validPreparationsIndexManager.Index(ctx, msg.ValidPreparation.ID, msg.ValidPreparation); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the valid preparation")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "validPreparationUpdated",
				ValidPreparation:          msg.ValidPreparation,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.ValidIngredientPreparationDataType:
		if err := w.dataManager.UpdateValidIngredientPreparation(ctx, msg.ValidIngredientPreparation); err != nil {
			return observability.PrepareError(err, logger, span, "creating valid ingredient preparation")
		}

		if err := w.validIngredientPreparationsIndexManager.Index(ctx, msg.ValidIngredientPreparation.ID, msg.ValidIngredientPreparation); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the valid ingredient preparation")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                   msg.DataType,
				MessageType:                "validIngredientPreparationUpdated",
				ValidIngredientPreparation: msg.ValidIngredientPreparation,
				AttributableToUserID:       msg.AttributableToUserID,
				AttributableToHouseholdID:  msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeDataType:
		if err := w.dataManager.UpdateRecipe(ctx, msg.Recipe); err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe")
		}

		if err := w.recipesIndexManager.Index(ctx, msg.Recipe.ID, msg.Recipe); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeUpdated",
				Recipe:                    msg.Recipe,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeStepDataType:
		if err := w.dataManager.UpdateRecipeStep(ctx, msg.RecipeStep); err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe step")
		}

		if err := w.recipeStepsIndexManager.Index(ctx, msg.RecipeStep.ID, msg.RecipeStep); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe step")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeStepUpdated",
				RecipeStep:                msg.RecipeStep,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeStepInstrumentDataType:
		if err := w.dataManager.UpdateRecipeStepInstrument(ctx, msg.RecipeStepInstrument); err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe step instrument")
		}

		if err := w.recipeStepInstrumentsIndexManager.Index(ctx, msg.RecipeStepInstrument.ID, msg.RecipeStepInstrument); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe step instrument")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeStepInstrumentUpdated",
				RecipeStepInstrument:      msg.RecipeStepInstrument,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeStepIngredientDataType:
		if err := w.dataManager.UpdateRecipeStepIngredient(ctx, msg.RecipeStepIngredient); err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe step ingredient")
		}

		if err := w.recipeStepIngredientsIndexManager.Index(ctx, msg.RecipeStepIngredient.ID, msg.RecipeStepIngredient); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe step ingredient")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeStepIngredientUpdated",
				RecipeStepIngredient:      msg.RecipeStepIngredient,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeStepProductDataType:
		if err := w.dataManager.UpdateRecipeStepProduct(ctx, msg.RecipeStepProduct); err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe step product")
		}

		if err := w.recipeStepProductsIndexManager.Index(ctx, msg.RecipeStepProduct.ID, msg.RecipeStepProduct); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe step product")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeStepProductUpdated",
				RecipeStepProduct:         msg.RecipeStepProduct,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.MealPlanDataType:
		if err := w.dataManager.UpdateMealPlan(ctx, msg.MealPlan); err != nil {
			return observability.PrepareError(err, logger, span, "creating meal plan")
		}

		if err := w.mealPlansIndexManager.Index(ctx, msg.MealPlan.ID, msg.MealPlan); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the meal plan")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "mealPlanUpdated",
				MealPlan:                  msg.MealPlan,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.MealPlanOptionDataType:
		if err := w.dataManager.UpdateMealPlanOption(ctx, msg.MealPlanOption); err != nil {
			return observability.PrepareError(err, logger, span, "creating meal plan option")
		}

		if err := w.mealPlanOptionsIndexManager.Index(ctx, msg.MealPlanOption.ID, msg.MealPlanOption); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the meal plan option")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "mealPlanOptionUpdated",
				MealPlanOption:            msg.MealPlanOption,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.MealPlanOptionVoteDataType:
		if err := w.dataManager.UpdateMealPlanOptionVote(ctx, msg.MealPlanOptionVote); err != nil {
			return observability.PrepareError(err, logger, span, "creating meal plan option vote")
		}

		if err := w.mealPlanOptionVotesIndexManager.Index(ctx, msg.MealPlanOptionVote.ID, msg.MealPlanOptionVote); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the meal plan option vote")
		}

		if w.postUpdatesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "mealPlanOptionVoteUpdated",
				MealPlanOptionVote:        msg.MealPlanOptionVote,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postUpdatesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.UserMembershipDataType,
		types.WebhookDataType,
		types.HouseholdInvitationDataType:
		break
	}

	return nil
}

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

// PreWritesWorker writes data from the pending writes topic to the database.
type PreWritesWorker struct {
	logger                                  logging.Logger
	tracer                                  tracing.Tracer
	encoder                                 encoding.ClientEncoder
	postWritesPublisher                     publishers.Publisher
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

// ProvidePreWritesWorker provides a PreWritesWorker.
func ProvidePreWritesWorker(
	ctx context.Context,
	logger logging.Logger,
	client *http.Client,
	dataManager database.DataManager,
	postWritesPublisher publishers.Publisher,
	searchIndexLocation search.IndexPath,
	searchIndexProvider search.IndexManagerProvider,
) (*PreWritesWorker, error) {
	const name = "pre_writes"

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

	w := &PreWritesWorker{
		logger:                                  logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                                  tracing.NewTracer(name),
		encoder:                                 encoding.ProvideClientEncoder(logger, encoding.ContentTypeJSON),
		postWritesPublisher:                     postWritesPublisher,
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

// HandleMessage handles a pending write.
func (w *PreWritesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	w.logger.Debug("message received")
	var msg *types.PreWriteMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	logger := w.logger.WithValue("data_type", msg.DataType)

	logger.Debug("message read successfully")

	switch msg.DataType {
	case types.ValidInstrumentDataType:
		validInstrument, err := w.dataManager.CreateValidInstrument(ctx, msg.ValidInstrument)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating valid instrument")
		}

		if err = w.validInstrumentsIndexManager.Index(ctx, validInstrument.ID, validInstrument); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the valid instrument")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "validInstrumentCreated",
				ValidInstrument:           validInstrument,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.ValidIngredientDataType:
		validIngredient, err := w.dataManager.CreateValidIngredient(ctx, msg.ValidIngredient)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating valid ingredient")
		}

		if err = w.validIngredientsIndexManager.Index(ctx, validIngredient.ID, validIngredient); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the valid ingredient")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "validIngredientCreated",
				ValidIngredient:           validIngredient,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.ValidPreparationDataType:
		validPreparation, err := w.dataManager.CreateValidPreparation(ctx, msg.ValidPreparation)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating valid preparation")
		}

		if err = w.validPreparationsIndexManager.Index(ctx, validPreparation.ID, validPreparation); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the valid preparation")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "validPreparationCreated",
				ValidPreparation:          validPreparation,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.ValidIngredientPreparationDataType:
		validIngredientPreparation, err := w.dataManager.CreateValidIngredientPreparation(ctx, msg.ValidIngredientPreparation)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating valid ingredient preparation")
		}

		if err = w.validIngredientPreparationsIndexManager.Index(ctx, validIngredientPreparation.ID, validIngredientPreparation); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the valid ingredient preparation")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                   msg.DataType,
				MessageType:                "validIngredientPreparationCreated",
				ValidIngredientPreparation: validIngredientPreparation,
				AttributableToUserID:       msg.AttributableToUserID,
				AttributableToHouseholdID:  msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.RecipeDataType:
		recipe, err := w.dataManager.CreateRecipe(ctx, msg.Recipe)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe")
		}

		if err = w.recipesIndexManager.Index(ctx, recipe.ID, recipe); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeCreated",
				Recipe:                    recipe,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.RecipeStepDataType:
		recipeStep, err := w.dataManager.CreateRecipeStep(ctx, msg.RecipeStep)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe step")
		}

		if err = w.recipeStepsIndexManager.Index(ctx, recipeStep.ID, recipeStep); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe step")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeStepCreated",
				RecipeStep:                recipeStep,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.RecipeStepInstrumentDataType:
		recipeStepInstrument, err := w.dataManager.CreateRecipeStepInstrument(ctx, msg.RecipeStepInstrument)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe step instrument")
		}

		if err = w.recipeStepInstrumentsIndexManager.Index(ctx, recipeStepInstrument.ID, recipeStepInstrument); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe step instrument")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeStepInstrumentCreated",
				RecipeStepInstrument:      recipeStepInstrument,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.RecipeStepIngredientDataType:
		recipeStepIngredient, err := w.dataManager.CreateRecipeStepIngredient(ctx, msg.RecipeStepIngredient)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe step ingredient")
		}

		if err = w.recipeStepIngredientsIndexManager.Index(ctx, recipeStepIngredient.ID, recipeStepIngredient); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe step ingredient")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeStepIngredientCreated",
				RecipeStepIngredient:      recipeStepIngredient,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.RecipeStepProductDataType:
		recipeStepProduct, err := w.dataManager.CreateRecipeStepProduct(ctx, msg.RecipeStepProduct)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating recipe step product")
		}

		if err = w.recipeStepProductsIndexManager.Index(ctx, recipeStepProduct.ID, recipeStepProduct); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the recipe step product")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "recipeStepProductCreated",
				RecipeStepProduct:         recipeStepProduct,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.MealPlanDataType:
		mealPlan, err := w.dataManager.CreateMealPlan(ctx, msg.MealPlan)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating meal plan")
		}

		if err = w.mealPlansIndexManager.Index(ctx, mealPlan.ID, mealPlan); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the meal plan")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "mealPlanCreated",
				MealPlan:                  mealPlan,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.MealPlanOptionDataType:
		mealPlanOption, err := w.dataManager.CreateMealPlanOption(ctx, msg.MealPlanOption)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating meal plan option")
		}

		if err = w.mealPlanOptionsIndexManager.Index(ctx, mealPlanOption.ID, mealPlanOption); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the meal plan option")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "mealPlanOptionCreated",
				MealPlanOption:            mealPlanOption,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.MealPlanOptionVoteDataType:
		mealPlanOptionVote, err := w.dataManager.CreateMealPlanOptionVote(ctx, msg.MealPlanOptionVote)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating meal plan option vote")
		}

		if err = w.mealPlanOptionVotesIndexManager.Index(ctx, mealPlanOptionVote.ID, mealPlanOptionVote); err != nil {
			return observability.PrepareError(err, logger, span, "indexing the meal plan option vote")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "mealPlanOptionVoteCreated",
				MealPlanOptionVote:        mealPlanOptionVote,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
			}
		}
	case types.WebhookDataType:
		webhook, err := w.dataManager.CreateWebhook(ctx, msg.Webhook)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating webhook")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "webhookCreated",
				Webhook:                   webhook,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.HouseholdInvitationDataType:
		householdInvitation, err := w.dataManager.CreateHouseholdInvitation(ctx, msg.HouseholdInvitation)
		if err != nil {
			return observability.PrepareError(err, logger, span, "creating user membership")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "userMembershipCreated",
				HouseholdInvitation:       householdInvitation,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}
			if err = w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.UserMembershipDataType:
		if err := w.dataManager.AddUserToHousehold(ctx, msg.UserMembership); err != nil {
			return observability.PrepareError(err, logger, span, "creating user membership")
		}

		if w.postWritesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				MessageType:               "userMembershipCreated",
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}
			if err := w.postWritesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	default:
		return observability.PrepareError(fmt.Errorf("invalid message type: %q", msg.DataType), logger, span, "handling message")
	}

	return nil
}

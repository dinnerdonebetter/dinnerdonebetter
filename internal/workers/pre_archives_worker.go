package workers

import (
	"context"
	"fmt"
	"net/http"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	publishers "gitlab.com/prixfixe/prixfixe/internal/messagequeue/publishers"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// PreArchivesWorker archives data from the pending archives topic to the database.
type PreArchivesWorker struct {
	logger                                  logging.Logger
	tracer                                  tracing.Tracer
	encoder                                 encoding.ClientEncoder
	postArchivesPublisher                   publishers.Publisher
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

// ProvidePreArchivesWorker provides a PreArchivesWorker.
func ProvidePreArchivesWorker(
	ctx context.Context,
	logger logging.Logger,
	client *http.Client,
	dataManager database.DataManager,
	postArchivesPublisher publishers.Publisher,
	searchIndexLocation search.IndexPath,
	searchIndexProvider search.IndexManagerProvider,
) (*PreArchivesWorker, error) {
	const name = "pre_archives"

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

	w := &PreArchivesWorker{
		logger:                                  logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                                  tracing.NewTracer(name),
		encoder:                                 encoding.ProvideClientEncoder(logger, encoding.ContentTypeJSON),
		postArchivesPublisher:                   postArchivesPublisher,
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

// HandleMessage handles a pending archive.
func (w *PreArchivesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var msg *types.PreArchiveMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	logger := w.logger.WithValue("data_type", msg.DataType)

	logger.Debug("message read")

	switch msg.DataType {
	case types.ValidInstrumentDataType:
		if err := w.dataManager.ArchiveValidInstrument(ctx, msg.ValidInstrumentID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving valid instrument")
		}

		if err := w.validInstrumentsIndexManager.Delete(ctx, msg.ValidInstrumentID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing valid instrument from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.ValidIngredientDataType:
		if err := w.dataManager.ArchiveValidIngredient(ctx, msg.ValidIngredientID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving valid ingredient")
		}

		if err := w.validIngredientsIndexManager.Delete(ctx, msg.ValidIngredientID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing valid ingredient from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.ValidPreparationDataType:
		if err := w.dataManager.ArchiveValidPreparation(ctx, msg.ValidPreparationID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving valid preparation")
		}

		if err := w.validPreparationsIndexManager.Delete(ctx, msg.ValidPreparationID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing valid preparation from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.ValidIngredientPreparationDataType:
		if err := w.dataManager.ArchiveValidIngredientPreparation(ctx, msg.ValidIngredientPreparationID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving valid ingredient preparation")
		}

		if err := w.validIngredientPreparationsIndexManager.Delete(ctx, msg.ValidIngredientPreparationID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing valid ingredient preparation from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeDataType:
		if err := w.dataManager.ArchiveRecipe(ctx, msg.RecipeID, msg.AttributableToHouseholdID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving recipe")
		}

		if err := w.recipesIndexManager.Delete(ctx, msg.RecipeID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing recipe from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeStepDataType:
		if err := w.dataManager.ArchiveRecipeStep(ctx, msg.RecipeID, msg.RecipeStepID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving recipe step")
		}

		if err := w.recipeStepsIndexManager.Delete(ctx, msg.RecipeStepID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing recipe step from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeStepInstrumentDataType:
		if err := w.dataManager.ArchiveRecipeStepInstrument(ctx, msg.RecipeStepID, msg.RecipeStepInstrumentID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving recipe step instrument")
		}

		if err := w.recipeStepInstrumentsIndexManager.Delete(ctx, msg.RecipeStepInstrumentID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing recipe step instrument from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeStepIngredientDataType:
		if err := w.dataManager.ArchiveRecipeStepIngredient(ctx, msg.RecipeStepID, msg.RecipeStepIngredientID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving recipe step ingredient")
		}

		if err := w.recipeStepIngredientsIndexManager.Delete(ctx, msg.RecipeStepIngredientID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing recipe step ingredient from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.RecipeStepProductDataType:
		if err := w.dataManager.ArchiveRecipeStepProduct(ctx, msg.RecipeStepID, msg.RecipeStepProductID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving recipe step product")
		}

		if err := w.recipeStepProductsIndexManager.Delete(ctx, msg.RecipeStepProductID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing recipe step product from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.MealPlanDataType:
		if err := w.dataManager.ArchiveMealPlan(ctx, msg.MealPlanID, msg.AttributableToHouseholdID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving meal plan")
		}

		if err := w.mealPlansIndexManager.Delete(ctx, msg.MealPlanID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing meal plan from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.MealPlanOptionDataType:
		if err := w.dataManager.ArchiveMealPlanOption(ctx, msg.MealPlanOptionID, msg.AttributableToHouseholdID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving meal plan option")
		}

		if err := w.mealPlanOptionsIndexManager.Delete(ctx, msg.MealPlanOptionID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing meal plan option from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.MealPlanOptionVoteDataType:
		if err := w.dataManager.ArchiveMealPlanOptionVote(ctx, msg.MealPlanOptionVoteID, msg.AttributableToHouseholdID); err != nil {
			return observability.PrepareError(err, w.logger, span, "archiving meal plan option vote")
		}

		if err := w.mealPlanOptionVotesIndexManager.Delete(ctx, msg.MealPlanOptionVoteID); err != nil {
			return observability.PrepareError(err, w.logger, span, "removing meal plan option vote from index")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.WebhookDataType:
		if err := w.dataManager.ArchiveWebhook(ctx, msg.WebhookID, msg.AttributableToHouseholdID); err != nil {
			return observability.PrepareError(err, w.logger, span, "creating webhook")
		}

		if w.postArchivesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:                  msg.DataType,
				AttributableToUserID:      msg.AttributableToUserID,
				AttributableToHouseholdID: msg.AttributableToHouseholdID,
			}

			if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
				return observability.PrepareError(err, logger, span, "publishing data change message")
			}
		}
	case types.UserMembershipDataType:
		break
	}

	return nil
}

package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (w *WritesWorker) createRecipeStepIngredient(ctx context.Context, msg *types.PreWriteMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	recipeStepIngredient, err := w.dataManager.CreateRecipeStepIngredient(ctx, msg.RecipeStepIngredient)
	if err != nil {
		return observability.PrepareError(err, logger, span, "creating recipe step ingredient")
	}

	if w.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "recipeStepIngredientCreated",
			RecipeStepIngredient:      recipeStepIngredient,
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err = w.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
		}
	}

	return nil
}

func (w *UpdatesWorker) updateRecipeStepIngredient(ctx context.Context, msg *types.PreUpdateMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.UpdateRecipeStepIngredient(ctx, msg.RecipeStepIngredient); err != nil {
		return observability.PrepareError(err, logger, span, "creating recipe step ingredient")
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

	return nil
}

func (w *ArchivesWorker) archiveRecipeStepIngredient(ctx context.Context, msg *types.PreArchiveMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.ArchiveRecipeStepIngredient(ctx, msg.RecipeStepID, msg.RecipeStepIngredientID); err != nil {
		return observability.PrepareError(err, w.logger, span, "archiving recipe step ingredient")
	}

	if w.postArchivesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "recipeStepIngredientArchived",
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

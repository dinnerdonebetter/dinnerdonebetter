package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (w *WritesWorker) createRecipeStepInstrument(ctx context.Context, msg *types.PreWriteMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	recipeStepInstrument, err := w.dataManager.CreateRecipeStepInstrument(ctx, msg.RecipeStepInstrument)
	if err != nil {
		return observability.PrepareError(err, logger, span, "creating recipe step ingredient")
	}

	if w.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "recipeStepInstrumentCreated",
			RecipeStepInstrument:      recipeStepInstrument,
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err = w.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
		}
	}

	return nil
}

func (w *UpdatesWorker) updateRecipeStepInstrument(ctx context.Context, msg *types.PreUpdateMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.UpdateRecipeStepInstrument(ctx, msg.RecipeStepInstrument); err != nil {
		return observability.PrepareError(err, logger, span, "creating recipe step ingredient")
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

	return nil
}

func (w *ArchivesWorker) archiveRecipeStepInstrument(ctx context.Context, msg *types.PreArchiveMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.ArchiveRecipeStepInstrument(ctx, msg.RecipeStepID, msg.RecipeStepInstrumentID); err != nil {
		return observability.PrepareError(err, w.logger, span, "archiving recipe step ingredient")
	}

	if w.postArchivesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "recipeStepInstrumentArchived",
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

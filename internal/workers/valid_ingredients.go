package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (w *WritesWorker) createValidIngredient(ctx context.Context, msg *types.PreWriteMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

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

	return nil
}

func (w *UpdatesWorker) updateValidIngredient(ctx context.Context, msg *types.PreUpdateMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

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

	return nil
}

func (w *ArchivesWorker) archiveValidIngredient(ctx context.Context, msg *types.PreArchiveMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.ArchiveValidIngredient(ctx, msg.ValidIngredientID); err != nil {
		return observability.PrepareError(err, w.logger, span, "archiving valid ingredient")
	}

	if err := w.validIngredientsIndexManager.Delete(ctx, msg.ValidIngredientID); err != nil {
		return observability.PrepareError(err, w.logger, span, "removing valid ingredient from index")
	}

	if w.postArchivesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "validIngredientArchived",
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

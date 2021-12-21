package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (w *WritesWorker) createValidIngredientPreparation(ctx context.Context, msg *types.PreWriteMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	validIngredientPreparation, err := w.dataManager.CreateValidIngredientPreparation(ctx, msg.ValidIngredientPreparation)
	if err != nil {
		return observability.PrepareError(err, logger, span, "creating valid ingredient preparation")
	}

	if w.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                   msg.DataType,
			MessageType:                "validIngredientPreparationCreated",
			ValidIngredientPreparation: validIngredientPreparation,
			AttributableToUserID:       msg.AttributableToUserID,
			AttributableToHouseholdID:  msg.AttributableToHouseholdID,
		}

		if err = w.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
		}
	}

	return nil
}

func (w *UpdatesWorker) updateValidIngredientPreparation(ctx context.Context, msg *types.PreUpdateMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.UpdateValidIngredientPreparation(ctx, msg.ValidIngredientPreparation); err != nil {
		return observability.PrepareError(err, logger, span, "creating valid ingredient preparation")
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

	return nil
}

func (w *ArchivesWorker) archiveValidIngredientPreparation(ctx context.Context, msg *types.PreArchiveMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.ArchiveValidIngredientPreparation(ctx, msg.ValidIngredientPreparationID); err != nil {
		return observability.PrepareError(err, w.logger, span, "archiving valid ingredient preparation")
	}

	if w.postArchivesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "validIngredientPreparationArchived",
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

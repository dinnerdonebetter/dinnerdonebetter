package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (w *WritesWorker) createValidPreparation(ctx context.Context, msg *types.PreWriteMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

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

	return nil
}

func (w *UpdatesWorker) updateValidPreparation(ctx context.Context, msg *types.PreUpdateMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

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

	return nil
}

func (w *ArchivesWorker) archiveValidPreparation(ctx context.Context, msg *types.PreArchiveMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.ArchiveValidPreparation(ctx, msg.ValidPreparationID); err != nil {
		return observability.PrepareError(err, w.logger, span, "archiving valid preparation")
	}

	if err := w.validPreparationsIndexManager.Delete(ctx, msg.ValidPreparationID); err != nil {
		return observability.PrepareError(err, w.logger, span, "removing valid preparation from index")
	}

	if w.postArchivesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "validPreparationArchived",
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

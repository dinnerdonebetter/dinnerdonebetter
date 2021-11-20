package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (w *WritesWorker) createValidInstrument(ctx context.Context, msg *types.PreWriteMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	validInstrument, err := w.dataManager.CreateValidInstrument(ctx, msg.ValidInstrument)
	if err != nil {
		return observability.PrepareError(err, logger, span, "creating valid instrument")
	}

	if err = w.validInstrumentsIndexManager.Index(ctx, validInstrument.ID, validInstrument); err != nil {
		return observability.PrepareError(err, logger, span, "indexing the valid instrument")
	}

	if w.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "validInstrumentCreated",
			ValidInstrument:           validInstrument,
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err = w.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing to post-writes topic")
		}
	}

	return nil
}

func (w *UpdatesWorker) updateValidInstrument(ctx context.Context, msg *types.PreUpdateMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

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

	return nil
}

func (w *ArchivesWorker) archiveValidInstrument(ctx context.Context, msg *types.PreArchiveMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.ArchiveValidInstrument(ctx, msg.ValidInstrumentID); err != nil {
		return observability.PrepareError(err, w.logger, span, "archiving valid instrument")
	}

	if err := w.validInstrumentsIndexManager.Delete(ctx, msg.ValidInstrumentID); err != nil {
		return observability.PrepareError(err, w.logger, span, "removing valid instrument from index")
	}

	if w.postArchivesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "validInstrumentArchived",
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

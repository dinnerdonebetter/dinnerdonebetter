package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (w *WritesWorker) createMealPlanOptionVote(ctx context.Context, msg *types.PreWriteMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	mealPlanOptionVote, err := w.dataManager.CreateMealPlanOptionVote(ctx, msg.MealPlanOptionVote)
	if err != nil {
		return observability.PrepareError(err, logger, span, "creating meal plan option vote")
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

	// have all votes been received for an option? if so, finalize it
	finalized, err := w.dataManager.FinalizeMealPlanOption(ctx, msg.MealPlanID, mealPlanOptionVote.BelongsToMealPlanOption)
	if err != nil {
		return observability.PrepareError(err, logger, span, "finalizing meal plan option")
	}

	// fire event

	// have all options for the meal plan been selected? if so, finalize the meal plan and fire event
	if finalized {
		logger.Debug("meal plan finalized")
	}

	return nil
}

func (w *UpdatesWorker) updateMealPlanOptionVote(ctx context.Context, msg *types.PreUpdateMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.UpdateMealPlanOptionVote(ctx, msg.MealPlanOptionVote); err != nil {
		return observability.PrepareError(err, logger, span, "creating meal plan option vote")
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

	return nil
}

func (w *ArchivesWorker) archiveMealPlanOptionVote(ctx context.Context, msg *types.PreArchiveMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.ArchiveMealPlanOptionVote(ctx, msg.MealPlanOptionID, msg.MealPlanOptionVoteID); err != nil {
		return observability.PrepareError(err, w.logger, span, "archiving meal plan option vote")
	}

	if w.postArchivesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "mealPlanOptionVoteArchived",
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.postArchivesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

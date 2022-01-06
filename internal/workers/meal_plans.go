package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (w *WritesWorker) createMealPlan(ctx context.Context, msg *types.PreWriteMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	mealPlan, err := w.dataManager.CreateMealPlan(ctx, msg.MealPlan)
	if err != nil {
		return observability.PrepareError(err, logger, span, "creating meal plan")
	}

	if w.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "mealPlanCreated",
			MealPlan:                  mealPlan,
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err = w.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing to data changes topic")
		}
	}

	return nil
}

func (w *UpdatesWorker) updateMealPlan(ctx context.Context, msg *types.PreUpdateMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.UpdateMealPlan(ctx, msg.MealPlan); err != nil {
		return observability.PrepareError(err, logger, span, "creating meal plan")
	}

	if w.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "mealPlanUpdated",
			MealPlan:                  msg.MealPlan,
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

func (w *ArchivesWorker) archiveMealPlan(ctx context.Context, msg *types.PreArchiveMessage) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("data_type", msg.DataType)

	if err := w.dataManager.ArchiveMealPlan(ctx, msg.MealPlanID, msg.AttributableToHouseholdID); err != nil {
		return observability.PrepareError(err, w.logger, span, "archiving meal plan")
	}

	if w.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  msg.DataType,
			MessageType:               "mealPlanArchived",
			AttributableToUserID:      msg.AttributableToUserID,
			AttributableToHouseholdID: msg.AttributableToHouseholdID,
		}

		if err := w.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			return observability.PrepareError(err, logger, span, "publishing data change message")
		}
	}

	return nil
}

func (w *ChoresWorker) finalizeExpiredMealPlans(ctx context.Context, msg *types.ChoreMessage) error {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValues(map[string]interface{}{
		"chore.type":        msg.ChoreType,
		keys.MealPlanIDKey:  msg.MealPlanID,
		keys.HouseholdIDKey: msg.AttributableToHouseholdID,
	})

	logger.Debug("finalize meal plan chore invoked")

	changed, err := w.dataManager.FinalizeMealPlanWithExpiredVotingPeriod(ctx, msg.MealPlanID, msg.AttributableToHouseholdID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "finalizing meal plan")
	}

	if !changed {
		logger.Error(nil, "meal plan was not changed")
	}

	return nil
}

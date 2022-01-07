package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/pkg/types"
)

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

package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
)

func (w *ChoresWorker) finalizeExpiredMealPlans(ctx context.Context) error {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger

	mealPlans, fetchMealPlansErr := w.dataManager.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
	if fetchMealPlansErr != nil {
		logger.Fatal(fetchMealPlansErr)
	}

	for _, mealPlan := range mealPlans {
		changed, err := w.dataManager.FinalizeMealPlanWithExpiredVotingPeriod(ctx, mealPlan.ID, mealPlan.BelongsToHousehold)
		if err != nil {
			return observability.PrepareError(err, logger, span, "finalizing meal plan")
		}

		if !changed {
			logger.Debug("meal plan was not changed")
		}
	}

	return nil
}

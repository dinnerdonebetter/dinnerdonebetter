package workers

import (
	"context"

	"github.com/prixfixeco/backend/internal/analytics"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

// MealPlanTallyScheduler finalizes meal plans.
type MealPlanTallyScheduler struct {
	logger                 logging.Logger
	tracer                 tracing.Tracer
	encoder                encoding.ClientEncoder
	dataManager            database.DataManager
	tallyQueuePublisher    messagequeue.Publisher
	emailSender            email.Emailer
	analyticsEventReporter analytics.EventReporter
}

// ProvideMealPlanTallyScheduler provides a MealPlanTallyScheduler.
func ProvideMealPlanTallyScheduler(
	logger logging.Logger,
	dataManager database.DataManager,
	tallyQueuePublisher messagequeue.Publisher,
	emailSender email.Emailer,
	analyticsEventReporter analytics.EventReporter,
	tracerProvider tracing.TracerProvider,
) *MealPlanTallyScheduler {
	n := "meal_plan_tally_queueing_worker"

	return &MealPlanTallyScheduler{
		logger:                 logging.EnsureLogger(logger).WithName(n),
		tracer:                 tracing.NewTracer(tracerProvider.Tracer(n)),
		encoder:                encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:            dataManager,
		tallyQueuePublisher:    tallyQueuePublisher,
		emailSender:            emailSender,
		analyticsEventReporter: analyticsEventReporter,
	}
}

// ScheduleMealPlanTallies handles a message ordering the finalization of expired meal plans.
func (w *MealPlanTallyScheduler) ScheduleMealPlanTallies(ctx context.Context, _ []byte) error {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	mealPlans, fetchMealPlansErr := w.dataManager.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
	if fetchMealPlansErr != nil {
		return observability.PrepareAndLogError(fetchMealPlansErr, logger, span, "fetching unfinalized and expired meal plan")
	}

	logger.WithValue("quantity", len(mealPlans)).Info("finalizing expired meal plans")

	for _, mealPlan := range mealPlans {
		if err := w.tallyQueuePublisher.Publish(ctx, &MealPlanTallyRequest{
			MealPlanID:  mealPlan.ID,
			HouseholdID: mealPlan.BelongsToHousehold,
		}); err != nil {
			observability.AcknowledgeError(err, logger, span, "queueing tally request for meal plan")
		}
	}

	return nil
}

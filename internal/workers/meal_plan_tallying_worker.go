package workers

import (
	"context"
	"fmt"

	"github.com/prixfixeco/backend/internal/analytics"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

const (
	mealPlanTallyWorkerName = "meal_plan_tallier"
)

type (
	MealPlanTallyRequest struct {
		MealPlanID  string `json:"mealPlanID"`
		HouseholdID string `json:"householdID"`
	}

	// MealPlanTallyingWorker finalizes meal plans.
	MealPlanTallyingWorker struct {
		logger                 logging.Logger
		tracer                 tracing.Tracer
		encoder                encoding.ClientEncoder
		dataManager            database.DataManager
		postUpdatesPublisher   messagequeue.Publisher
		emailSender            email.Emailer
		analyticsEventReporter analytics.EventReporter
	}
)

// ProvideMealPlanTallyingWorker provides a MealPlanTallyingWorker.
func ProvideMealPlanTallyingWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	emailSender email.Emailer,
	analyticsEventReporter analytics.EventReporter,
	tracerProvider tracing.TracerProvider,
) *MealPlanTallyingWorker {
	return &MealPlanTallyingWorker{
		logger:                 logging.EnsureLogger(logger).WithName(mealPlanTallyWorkerName),
		tracer:                 tracing.NewTracer(tracerProvider.Tracer(mealPlanTallyWorkerName)),
		encoder:                encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:            dataManager,
		postUpdatesPublisher:   postUpdatesPublisher,
		emailSender:            emailSender,
		analyticsEventReporter: analyticsEventReporter,
	}
}

// TallyMealPlanVotes handles a message ordering the finalization of expired meal plans.
func (w *MealPlanTallyingWorker) TallyMealPlanVotes(ctx context.Context, req []byte) error {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var request *MealPlanTallyRequest
	if err := w.encoder.Unmarshal(ctx, req, &request); err != nil {
		return observability.PrepareError(err, span, "decoding request")
	}

	changed, err := w.dataManager.AttemptToFinalizeMealPlan(ctx, request.MealPlanID, request.HouseholdID)
	if err != nil {
		return observability.PrepareError(err, span, "finalizing meal plan")
	}

	if !changed {
		return fmt.Errorf("meal plan %s was not finalized", request.MealPlanID)
	}

	return nil
}

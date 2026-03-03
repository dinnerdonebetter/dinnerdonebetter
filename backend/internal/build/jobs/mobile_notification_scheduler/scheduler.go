package mobilenotificationscheduler

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/hashicorp/go-multierror"
)

const schedulerTracerName = "mobile_notification_scheduler"

// Scheduler publishes meal plan task notification requests to the mobile_notifications queue.
type Scheduler struct {
	logger                       logging.Logger
	tracer                       tracing.Tracer
	mealPlanningRepo             mealplanning.Repository
	mobileNotificationsPublisher messagequeue.Publisher
}

// NewScheduler creates a new mobile notification scheduler.
func NewScheduler(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	mealPlanningRepo mealplanning.Repository,
	mobileNotificationsPublisher messagequeue.Publisher,
) *Scheduler {
	return &Scheduler{
		logger:                       logging.EnsureLogger(logger).WithName(schedulerTracerName),
		tracer:                       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(schedulerTracerName)),
		mealPlanningRepo:             mealPlanningRepo,
		mobileNotificationsPublisher: mobileNotificationsPublisher,
	}
}

// ScheduleNotifications queries for meal plan tasks that need notifications and publishes them to the queue.
func (s *Scheduler) ScheduleNotifications(ctx context.Context) error {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	taskIDs, err := s.mealPlanningRepo.GetMealPlanTaskIDsThatNeedNotification(ctx)
	if err != nil {
		return observability.PrepareAndLogError(err, s.logger, span, "getting meal plan task IDs that need notification")
	}

	s.logger.WithValue("count", len(taskIDs)).Info("publishing mobile notification requests")

	errs := &multierror.Error{}
	for _, taskID := range taskIDs {
		req := &notifications.MealPlanTaskNotificationRequest{
			MealPlanTaskID: taskID,
		}
		if err = s.mobileNotificationsPublisher.Publish(ctx, req); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("publishing notification for task %s: %w", taskID, err))
		}
	}

	return errs.ErrorOrNil()
}

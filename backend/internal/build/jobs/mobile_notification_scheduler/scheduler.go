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

// ScheduleNotifications runs all notification schedulers. Each scheduler queries for items
// that need notifications and publishes them to the queue. Meal plan tasks are the first
// use case; additional notification types can be added by registering more handlers.
func (s *Scheduler) ScheduleNotifications(ctx context.Context) error {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	errs := &multierror.Error{}
	if err := s.scheduleMealPlanTaskNotifications(ctx); err != nil {
		errs = multierror.Append(errs, err)
	}
	// Future notification types: add more handler calls here.

	return errs.ErrorOrNil()
}

// scheduleMealPlanTaskNotifications queries for meal plan tasks that need notifications
// and publishes them to the mobile_notifications queue.
func (s *Scheduler) scheduleMealPlanTaskNotifications(ctx context.Context) error {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	taskIDs, err := s.mealPlanningRepo.GetMealPlanTaskIDsThatNeedNotification(ctx)
	if err != nil {
		return observability.PrepareAndLogError(err, s.logger, span, "getting meal plan task IDs that need notification")
	}

	if len(taskIDs) == 0 {
		return nil
	}

	s.logger.WithValue("count", len(taskIDs)).WithValue("type", "meal_plan_task").Info("publishing mobile notification requests")

	errs := &multierror.Error{}
	for _, taskID := range taskIDs {
		req := &notifications.MealPlanTaskNotificationRequest{
			MealPlanTaskID: taskID,
		}
		if err = s.mobileNotificationsPublisher.Publish(ctx, req); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("publishing meal plan task notification for %s: %w", taskID, err))
		}
	}

	return errs.ErrorOrNil()
}

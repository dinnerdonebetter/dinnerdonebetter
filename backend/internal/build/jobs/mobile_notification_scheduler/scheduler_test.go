package mobilenotificationscheduler

import (
	"testing"
	"time"

	identitymock "github.com/dinnerdonebetter/backend/internal/domain/identity/mock"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanningmock "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/mocks"
	msgqueuemock "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/reflection"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestScheduler_ScheduleNotifications_publishesMobileNotificationRequest(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	logger := logging.NewNoopLogger()
	tracerProvider := tracing.NewNoopTracerProvider()

	taskID := fakes.BuildFakeID()
	task := fakes.BuildFakeMealPlanTask()
	task.ID = taskID
	task.RecipePrepTask.Name = "Chop onions"
	assignedUser := fakes.BuildFakeID()
	task.AssignedToUser = &assignedUser

	mealPlanRepo := &mealplanningmock.Repository{}
	identityRepo := &identitymock.RepositoryMock{}
	publisher := &msgqueuemock.Publisher{}

	notificationCtx := &mealplanning.MealPlanTaskNotificationContext{
		PrepTaskName:        "Chop onions",
		CreationExplanation: "",
		MealName:            mealplanning.DinnerMealName,
		StartsAt:            time.Date(2025, 3, 3, 18, 0, 0, 0, time.UTC), // Monday
	}

	mealPlanRepo.On(reflection.GetMethodName(mealPlanRepo.GetMealPlanTaskIDsThatNeedNotification), mock.Anything).Return([]string{taskID}, nil).Once()
	mealPlanRepo.On(reflection.GetMethodName(mealPlanRepo.GetMealPlanTask), mock.Anything, taskID).Return(task, nil).Once()
	mealPlanRepo.On(reflection.GetMethodName(mealPlanRepo.GetMealPlanTaskNotificationContext), mock.Anything, taskID).Return(notificationCtx, nil).Once()
	// With AssignedToUser set, GetMealPlanTaskAccountID and GetUsersForAccount are not called

	var publishedPayload interface{}
	publisher.On(reflection.GetMethodName(publisher.Publish), mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		publishedPayload = args.Get(1)
	}).Return(nil).Once()

	scheduler := NewScheduler(logger, tracerProvider, mealPlanRepo, identityRepo, publisher)

	err := scheduler.ScheduleNotifications(ctx)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, mealPlanRepo, publisher)

	req, ok := publishedPayload.(*notifications.MobileNotificationRequest)
	require.True(t, ok, "expected MobileNotificationRequest to be published")
	assert.Equal(t, []string{assignedUser}, req.RecipientUserIDs)
	assert.Equal(t, "Meal plan task", req.Title)
	assert.Equal(t, "Chop onions for Dinner on Monday", req.Body)
	assert.NotNil(t, req.Context)
	assert.Equal(t, taskID, req.Context[notifications.MealPlanTaskIDContextKey])
}

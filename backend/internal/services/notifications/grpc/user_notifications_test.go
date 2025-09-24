package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	notificationsfakes "github.com/dinnerdonebetter/backend/internal/domain/notifications/fakes"
	notificationsmock "github.com/dinnerdonebetter/backend/internal/domain/notifications/mock"
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func buildTestService(t *testing.T) (*serviceImpl, *notificationsmock.Repository) {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	notificationsRepository := &notificationsmock.Repository{}

	service := &serviceImpl{
		tracer:                  tracer,
		logger:                  logger,
		notificationsRepository: notificationsRepository,
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: "test-account-id",
				Requester: sessions.RequesterInfo{
					UserID: "test-user-id",
				},
			}, nil
		},
	}

	return service, notificationsRepository
}

func buildTestServiceWithSessionError(t *testing.T) *serviceImpl {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())

	service := &serviceImpl{
		tracer: tracer,
		logger: logger,
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return nil, errors.New("session error")
		},
	}

	return service
}

func TestServiceImpl_GetUserNotification(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotification := notificationsfakes.BuildFakeUserNotification()
		userID := "test-user-id"

		service, notificationsRepository := buildTestService(t)

		request := &notificationssvc.GetUserNotificationRequest{
			UserNotificationID: userNotification.ID,
		}

		notificationsRepository.On("GetUserNotification", mock.MatchedBy(func(ctx context.Context) bool { return true }), userID, userNotification.ID).Return(userNotification, nil)

		response, err := service.GetUserNotification(ctx, request)

		require.NotNil(t, response)
		assert.NoError(t, err)
		assert.NotNil(t, response.ResponseDetails)
		assert.NotEmpty(t, response.ResponseDetails.TraceID)
		assert.NotNil(t, response.Result)
		assert.Equal(t, userNotification.ID, response.Result.ID)
		assert.Equal(t, userNotification.Content, response.Result.Content)

		mock.AssertExpectationsForObjects(t, notificationsRepository)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service := buildTestServiceWithSessionError(t)

		request := &notificationssvc.GetUserNotificationRequest{
			UserNotificationID: "test-notification-id",
		}

		response, err := service.GetUserNotification(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcStatus.Code())
	})

	t.Run("with repository error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotificationID := "test-notification-id"
		userID := "test-user-id"
		expectedError := errors.New("repository error")

		service, notificationsRepository := buildTestService(t)

		request := &notificationssvc.GetUserNotificationRequest{
			UserNotificationID: userNotificationID,
		}

		notificationsRepository.On("GetUserNotification", mock.MatchedBy(func(ctx context.Context) bool { return true }), userID, userNotificationID).Return((*notifications.UserNotification)(nil), expectedError)

		response, err := service.GetUserNotification(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, grpcStatus.Code())

		mock.AssertExpectationsForObjects(t, notificationsRepository)
	})
}

func TestServiceImpl_GetUserNotifications(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotifications := []*notifications.UserNotification{
			notificationsfakes.BuildFakeUserNotification(),
			notificationsfakes.BuildFakeUserNotification(),
		}
		userID := "test-user-id"

		service, notificationsRepository := buildTestService(t)

		request := &notificationssvc.GetUserNotificationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: func(x uint32) *uint32 { return &x }(20),
			},
		}

		expectedResult := &filtering.QueryFilteredResult[notifications.UserNotification]{
			Data: userNotifications,
			Pagination: filtering.Pagination{
				Page:          1,
				Limit:         20,
				FilteredCount: uint64(len(userNotifications)),
				TotalCount:    uint64(len(userNotifications)),
			},
		}

		notificationsRepository.On("GetUserNotifications", mock.MatchedBy(func(ctx context.Context) bool { return true }), userID, mock.AnythingOfType("*filtering.QueryFilter")).Return(expectedResult, nil)

		response, err := service.GetUserNotifications(ctx, request)

		require.NotNil(t, response)
		assert.NoError(t, err)
		assert.NotNil(t, response.ResponseDetails)
		assert.NotEmpty(t, response.ResponseDetails.TraceID)
		assert.NotNil(t, response.Results)
		assert.Len(t, response.Results, len(userNotifications))

		for i, result := range response.Results {
			assert.Equal(t, userNotifications[i].ID, result.ID)
			assert.Equal(t, userNotifications[i].Content, result.Content)
		}

		mock.AssertExpectationsForObjects(t, notificationsRepository)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service := buildTestServiceWithSessionError(t)

		request := &notificationssvc.GetUserNotificationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: func(x uint32) *uint32 { return &x }(20),
			},
		}

		response, err := service.GetUserNotifications(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcStatus.Code())
	})

	t.Run("with repository error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userID := "test-user-id"
		expectedError := errors.New("repository error")

		service, notificationsRepository := buildTestService(t)

		request := &notificationssvc.GetUserNotificationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: func(x uint32) *uint32 { return &x }(20),
			},
		}

		notificationsRepository.On("GetUserNotifications", mock.MatchedBy(func(ctx context.Context) bool { return true }), userID, mock.AnythingOfType("*filtering.QueryFilter")).Return((*filtering.QueryFilteredResult[notifications.UserNotification])(nil), expectedError)

		response, err := service.GetUserNotifications(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, grpcStatus.Code())

		mock.AssertExpectationsForObjects(t, notificationsRepository)
	})
}

func TestServiceImpl_UpdateUserNotification(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotification := notificationsfakes.BuildFakeUserNotification()
		userID := "test-user-id"
		newStatus := notifications.UserNotificationStatusTypeRead

		service, notificationsRepository := buildTestService(t)

		request := &notificationssvc.UpdateUserNotificationRequest{
			UserNotificationID: userNotification.ID,
			Input: &notificationssvc.UserNotificationUpdateRequestInput{
				Status: &newStatus,
			},
		}

		updatedNotification := &notifications.UserNotification{
			ID: userNotification.ID,
		}

		notificationsRepository.On("GetUserNotification", mock.MatchedBy(func(ctx context.Context) bool { return true }), userID, userNotification.ID).Return(userNotification, nil).Twice()
		notificationsRepository.On("UpdateUserNotification", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("*notifications.UserNotification")).Return(nil).Once()

		response, err := service.UpdateUserNotification(ctx, request)

		require.NotNil(t, response)
		assert.NoError(t, err)
		assert.NotNil(t, response.ResponseDetails)
		assert.NotEmpty(t, response.ResponseDetails.TraceID)
		assert.NotNil(t, response.Updated)
		assert.Equal(t, updatedNotification.ID, response.Updated.ID)

		mock.AssertExpectationsForObjects(t, notificationsRepository)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service := buildTestServiceWithSessionError(t)
		newStatus := notifications.UserNotificationStatusTypeRead

		request := &notificationssvc.UpdateUserNotificationRequest{
			UserNotificationID: "test-notification-id",
			Input: &notificationssvc.UserNotificationUpdateRequestInput{
				Status: &newStatus,
			},
		}

		response, err := service.UpdateUserNotification(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcStatus.Code())
	})

	t.Run("with get notification error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotificationID := "test-notification-id"
		userID := "test-user-id"
		expectedError := errors.New("repository error")
		newStatus := notifications.UserNotificationStatusTypeRead

		service, notificationsRepository := buildTestService(t)

		request := &notificationssvc.UpdateUserNotificationRequest{
			UserNotificationID: userNotificationID,
			Input: &notificationssvc.UserNotificationUpdateRequestInput{
				Status: &newStatus,
			},
		}

		notificationsRepository.On("GetUserNotification", mock.MatchedBy(func(ctx context.Context) bool { return true }), userID, userNotificationID).Return((*notifications.UserNotification)(nil), expectedError)

		response, err := service.UpdateUserNotification(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, grpcStatus.Code())

		mock.AssertExpectationsForObjects(t, notificationsRepository)
	})

	t.Run("with update notification error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotification := notificationsfakes.BuildFakeUserNotification()
		userID := "test-user-id"
		expectedError := errors.New("update error")
		newStatus := notifications.UserNotificationStatusTypeRead

		service, notificationsRepository := buildTestService(t)

		request := &notificationssvc.UpdateUserNotificationRequest{
			UserNotificationID: userNotification.ID,
			Input: &notificationssvc.UserNotificationUpdateRequestInput{
				Status: &newStatus,
			},
		}

		notificationsRepository.On("GetUserNotification", mock.MatchedBy(func(ctx context.Context) bool { return true }), userID, userNotification.ID).Return(userNotification, nil)
		notificationsRepository.On("UpdateUserNotification", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("*notifications.UserNotification")).Return(expectedError)

		response, err := service.UpdateUserNotification(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, grpcStatus.Code())

		mock.AssertExpectationsForObjects(t, notificationsRepository)
	})
}

package manager

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications/fakes"
	notificationsmock "github.com/dinnerdonebetter/backend/internal/domain/notifications/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/reflection"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildNotificationsManagerForTest(t *testing.T) *notificationsManager {
	t.Helper()

	ctx := context.Background()
	queueCfg := &msgconfig.QueuesConfig{
		DataChangesTopicName: t.Name(),
	}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On(reflection.GetMethodName(mpp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewNotificationsDataManager(
		ctx,
		tracing.NewNoopTracerProvider(),
		logging.NewNoopLogger(),
		&notificationsmock.Repository{},
		queueCfg,
		mpp,
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*notificationsManager)
}

func setupExpectationsForNotificationsManager(
	manager *notificationsManager,
	repoSetupFunc func(repo *notificationsmock.Repository),
	eventTypeMaps ...map[string][]string,
) []any {
	repo := &notificationsmock.Repository{}
	if repoSetupFunc != nil {
		repoSetupFunc(repo)
	}
	manager.repo = repo

	mp := &mockpublishers.Publisher{}
	for _, eventTypeMap := range eventTypeMaps {
		for eventType, payload := range eventTypeMap {
			mp.On(reflection.GetMethodName(mp.PublishAsync), testutils.ContextMatcher, eventMatches(eventType, payload)).Return()
		}
	}
	manager.dataChangesPublisher = mp

	return []any{repo, mp}
}

func TestNotificationsManager_CreateUserNotification(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		nm := buildNotificationsManagerForTest(t)

		expected := fakes.BuildFakeUserNotification()
		input := converters.ConvertUserNotificationToUserNotificationDatabaseCreationInput(expected)

		expectations := setupExpectationsForNotificationsManager(
			nm,
			func(repo *notificationsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.CreateUserNotification), testutils.ContextMatcher, testutils.MatchType[*notifications.UserNotificationDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				notifications.UserNotificationCreatedServiceEventType: {keys.UserNotificationIDKey},
			},
		)

		actual, err := nm.CreateUserNotification(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestNotificationsManager_UpdateUserNotification(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		nm := buildNotificationsManagerForTest(t)

		updated := fakes.BuildFakeUserNotification()
		updated.Status = notifications.UserNotificationStatusTypeRead

		expectations := setupExpectationsForNotificationsManager(
			nm,
			func(repo *notificationsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.UpdateUserNotification), testutils.ContextMatcher, updated).Return(nil)
			},
			map[string][]string{
				notifications.UserNotificationUpdatedServiceEventType: {keys.UserNotificationIDKey},
			},
		)

		err := nm.UpdateUserNotification(ctx, updated)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

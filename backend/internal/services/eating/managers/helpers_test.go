package managers

import (
	"context"
	"slices"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/events"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func eventMatches(eventType string, keys []string) any {
	return mock.MatchedBy(func(message *types.DataChangeMessage) bool {
		allContextKeys := []string{}
		for k := range message.Context {
			allContextKeys = append(allContextKeys, k)
		}

		slices.Sort(keys)
		slices.Sort(allContextKeys)
		allKeysMatch := slices.Equal(keys, allContextKeys)
		eventTypeMatches := message.EventType == eventType
		result := allKeysMatch && eventTypeMatches

		return result
	})
}

func setupExpectations(
	manager *mealPlanningManager,
	dbSetupFunc func(db *database.MockDatabase),
	eventTypeMaps ...map[string][]string,
) []any {
	db := database.NewMockDatabase()
	if dbSetupFunc != nil {
		dbSetupFunc(db)
	}
	manager.db = db

	mp := &mockpublishers.Publisher{}
	for _, eventTypeMap := range eventTypeMaps {
		for eventType, payload := range eventTypeMap {
			mp.On("PublishAsync", testutils.ContextMatcher, eventMatches(eventType, payload)).Return()
		}
	}
	manager.dataChangesPublisher = mp

	return []any{db, mp}
}

func buildMealPlanManagerForTest(t *testing.T) *mealPlanningManager {
	t.Helper()

	queueCfg := &msgconfig.QueuesConfig{
		DataChangesTopicName: t.Name(),
	}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On("ProvidePublisher", queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewMealPlanningManager(
		t.Context(),
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		database.NewMockDatabase(),
		queueCfg,
		mpp,
		&textsearchcfg.Config{},
		metrics.NewNoopMetricsProvider(),
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*mealPlanningManager)
}

func Test_buildDataChangeMessageFromContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		sessionContextData := &sessions.ContextData{
			Requester:         sessions.RequesterInfo{UserID: fakes.BuildFakeID()},
			ActiveHouseholdID: fakes.BuildFakeID(),
		}
		ctx = context.WithValue(ctx, sessions.SessionContextDataKey, sessionContextData)

		expected := &types.DataChangeMessage{
			EventType: events.MealCreated,
			Context: map[string]any{
				"things": "stuff",
			},
			UserID:      sessionContextData.Requester.UserID,
			HouseholdID: sessionContextData.ActiveHouseholdID,
		}

		actual := buildDataChangeMessageFromContext(ctx, logging.NewNoopLogger(), expected.EventType, expected.Context)

		assert.Equal(t, expected, actual)
	})
}

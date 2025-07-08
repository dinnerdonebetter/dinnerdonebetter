package managers

import (
	"context"
	"slices"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func eventMatches(eventType string, keys []string) any {
	return mock.MatchedBy(func(message *audit.DataChangeMessage) bool {
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

func Test_buildDataChangeMessageFromContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		sessionContextData := &sessions.ContextData{
			Requester:       sessions.RequesterInfo{UserID: fakes.BuildFakeID()},
			ActiveAccountID: fakes.BuildFakeID(),
		}
		ctx = context.WithValue(ctx, sessions.SessionContextDataKey, sessionContextData)

		expected := &audit.DataChangeMessage{
			EventType: mealplanning.MealCreatedServiceEventType,
			Context: map[string]any{
				"things": "stuff",
			},
			UserID:    sessionContextData.Requester.UserID,
			AccountID: sessionContextData.ActiveAccountID,
		}

		actual := audit.BuildDataChangeMessageFromContext(ctx, logging.NewNoopLogger(), expected.EventType, expected.Context)

		assert.Equal(t, expected, actual)
	})
}

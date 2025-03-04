package managers

import (
	"context"
	"slices"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/services/eating/events"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

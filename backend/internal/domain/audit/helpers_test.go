package audit

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

		expected := &DataChangeMessage{
			EventType: mealplanning.MealCreatedServiceEventType,
			Context: map[string]any{
				"things": "stuff",
			},
			UserID:    sessionContextData.Requester.UserID,
			AccountID: sessionContextData.ActiveAccountID,
		}

		actual := BuildDataChangeMessageFromContext(ctx, logging.NewNoopLogger(), expected.EventType, expected.Context)

		assert.Equal(t, expected, actual)
	})
}

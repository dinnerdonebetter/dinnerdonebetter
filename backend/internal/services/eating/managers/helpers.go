package managers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

func buildDataChangeMessageFromContext(ctx context.Context, logger logging.Logger, eventType string, metadata map[string]any) *types.DataChangeMessage {
	sessionContext, ok := ctx.Value(sessions.SessionContextDataKey).(*sessions.ContextData)
	if !ok {
		logger.WithValue("event_type", eventType).Info("failed to extract session data from context")
	}

	x := &types.DataChangeMessage{
		EventType: eventType,
		Context:   metadata,
	}

	if sessionContext != nil {
		x.UserID = sessionContext.Requester.UserID
		x.HouseholdID = sessionContext.ActiveHouseholdID
	}

	return x
}

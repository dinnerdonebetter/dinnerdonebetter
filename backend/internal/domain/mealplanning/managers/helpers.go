package managers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
)

func buildDataChangeMessageFromContext(ctx context.Context, logger logging.Logger, eventType string, metadata map[string]any) *audit.DataChangeMessage {
	sessionContext, ok := ctx.Value(sessions.SessionContextDataKey).(*sessions.ContextData)
	if !ok {
		logger.WithValue("event_type", eventType).Info("failed to extract session data from context")
	}

	x := &audit.DataChangeMessage{
		EventType: eventType,
		Context:   metadata,
	}

	if sessionContext != nil {
		x.UserID = sessionContext.Requester.UserID
		x.AccountID = sessionContext.ActiveAccountID
	}

	return x
}

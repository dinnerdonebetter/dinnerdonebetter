package analytics

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"
)

type (
	// EventReporter collects data about customers.
	EventReporter interface {
		Close()
		AddUser(ctx context.Context, userID string, properties map[string]any) error
		EventOccurred(ctx context.Context, event types.CustomerEventType, userID string, properties map[string]any) error
	}
)

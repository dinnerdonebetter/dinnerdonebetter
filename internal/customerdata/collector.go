package customerdata

import (
	"context"
	"io"

	"github.com/prixfixeco/api_server/pkg/types"
)

type (
	// Collector collects data about customers.
	Collector interface {
		io.Closer
		AddUser(ctx context.Context, userID string, properties map[string]interface{}) error
		EventOccurred(ctx context.Context, event types.CustomerEventType, userID string, properties map[string]interface{}) error
	}
)

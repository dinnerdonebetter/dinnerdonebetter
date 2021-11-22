package customerdata

import (
	"context"
	"io"
)

type (
	// Collector collects data about customers.
	Collector interface {
		io.Closer
		AddUser(ctx context.Context, userID string, properties map[string]interface{}) error
		EventOccurred(ctx context.Context, event, userID string, properties map[string]interface{}) error
	}
)

package customerdata

import (
	"context"
	"io"
)

type (
	// Collector collects data about customers.
	Collector interface {
		io.Closer
		Identify(ctx context.Context, userID string, properties map[string]interface{}) error
		Track(ctx context.Context, event, userID string, properties map[string]interface{}) error
	}
)

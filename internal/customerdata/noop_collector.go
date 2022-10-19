package customerdata

import (
	"context"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ Collector = (*noopCollector)(nil)

type (
	// noopCollector is a no-op Collector.
	noopCollector struct{}
)

// NewNoopCollector returns a new no-op Collector.
func NewNoopCollector() Collector {
	return &noopCollector{}
}

// Close does nothing.
func (c *noopCollector) Close() error {
	return nil
}

// AddUser does nothing.
func (c *noopCollector) AddUser(context.Context, string, map[string]interface{}) error {
	return nil
}

// EventOccurred does nothing.
func (c *noopCollector) EventOccurred(context.Context, types.CustomerEventType, string, map[string]interface{}) error {
	return nil
}

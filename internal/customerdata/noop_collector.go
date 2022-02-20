package customerdata

import (
	"context"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ Collector = (*NoopCollector)(nil)

type (
	// NoopCollector is a no-op customerdata.Collector.
	NoopCollector struct{}
)

// NewNoopCollector returns a new no-op Collector.
func NewNoopCollector() (*NoopCollector, error) {
	c := &NoopCollector{}

	return c, nil
}

// Close does nothing.
func (c *NoopCollector) Close() error {
	return nil
}

// AddUser does nothing.
func (c *NoopCollector) AddUser(context.Context, string, map[string]interface{}) error {
	return nil
}

// EventOccurred does nothing.
func (c *NoopCollector) EventOccurred(context.Context, types.CustomerEventType, string, map[string]interface{}) error {
	return nil
}

package customerdata

import (
	"context"
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

// Close wraps the internal client's Close method.
func (c *NoopCollector) Close() error {
	return nil
}

// Identify upsert's a user's identity.
func (c *NoopCollector) Identify(context.Context, string, map[string]interface{}) error {
	return nil
}

// Track associates events with a user.
func (c *NoopCollector) Track(context.Context, string, string, map[string]interface{}) error {
	return nil
}

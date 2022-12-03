package analytics

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ EventReporter = (*NoopEventReporter)(nil)

type (
	// NoopEventReporter is a no-op EventReporter.
	NoopEventReporter struct{}
)

// NewNoopCollector returns a new no-op EventReporter.
func NewNoopCollector() EventReporter {
	return &NoopEventReporter{}
}

// Close does nothing.
func (c *NoopEventReporter) Close() {}

// AddUser does nothing.
func (c *NoopEventReporter) AddUser(context.Context, string, map[string]any) error {
	return nil
}

// EventOccurred does nothing.
func (c *NoopEventReporter) EventOccurred(context.Context, types.CustomerEventType, string, map[string]any) error {
	return nil
}

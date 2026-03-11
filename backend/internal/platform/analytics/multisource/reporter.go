package multisource

import (
	"context"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/platform/analytics"
)

const name = "multisource_event_reporter"

// MultiSourceEventReporter delegates events to per-source EventReporters.
type MultiSourceEventReporter struct {
	reporters map[string]analytics.EventReporter
	mu        sync.RWMutex
}

// NewMultiSourceEventReporter returns a new MultiSourceEventReporter.
func NewMultiSourceEventReporter(reporters map[string]analytics.EventReporter) *MultiSourceEventReporter {
	if reporters == nil {
		reporters = make(map[string]analytics.EventReporter)
	}
	return &MultiSourceEventReporter{reporters: reporters}
}

// getReporter returns the reporter for the source, or Noop if unknown/missing.
func (m *MultiSourceEventReporter) getReporter(source string) analytics.EventReporter {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if r, ok := m.reporters[source]; ok && r != nil {
		return r
	}
	return analytics.NewNoopEventReporter()
}

// TrackEvent records an event for an identified user.
func (m *MultiSourceEventReporter) TrackEvent(ctx context.Context, source, event, userID string, properties map[string]any) error {
	return m.getReporter(source).EventOccurred(ctx, event, userID, properties)
}

// TrackAnonymousEvent records an event for an anonymous user.
func (m *MultiSourceEventReporter) TrackAnonymousEvent(ctx context.Context, source, event, anonymousID string, properties map[string]any) error {
	return m.getReporter(source).EventOccurredAnonymous(ctx, event, anonymousID, properties)
}

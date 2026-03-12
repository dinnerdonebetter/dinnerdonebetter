package multisource

import (
	"context"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/platform/analytics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const name = "multisource_event_reporter"

// MultiSourceEventReporter delegates events to per-source EventReporters.
type MultiSourceEventReporter struct {
	tracer    tracing.Tracer
	logger    logging.Logger
	reporters map[string]analytics.EventReporter
	mu        sync.RWMutex
}

// NewMultiSourceEventReporter returns a new MultiSourceEventReporter.
func NewMultiSourceEventReporter(
	reporters map[string]analytics.EventReporter,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
) *MultiSourceEventReporter {
	if reporters == nil {
		reporters = make(map[string]analytics.EventReporter)
	}
	return &MultiSourceEventReporter{
		reporters: reporters,
		logger:    logging.EnsureLogger(logger).WithName(name),
		tracer:    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(name)),
	}
}

// getReporter returns the reporter for the source, or Noop if unknown/missing.
func (m *MultiSourceEventReporter) getReporter(source string) analytics.EventReporter {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if r, ok := m.reporters[source]; ok && r != nil {
		return r
	}
	m.logger.WithValue("source", source).WithValue("known_sources", m.knownSources()).Info("no analytics reporter configured for source, using noop")
	return analytics.NewNoopEventReporter()
}

func (m *MultiSourceEventReporter) knownSources() []string {
	sources := make([]string, 0, len(m.reporters))
	for k := range m.reporters {
		sources = append(sources, k)
	}
	return sources
}

// TrackEvent records an event for an identified user.
func (m *MultiSourceEventReporter) TrackEvent(ctx context.Context, source, event, userID string, properties map[string]any) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.getReporter(source).EventOccurred(ctx, event, userID, properties)
}

// TrackAnonymousEvent records an event for an anonymous user.
func (m *MultiSourceEventReporter) TrackAnonymousEvent(ctx context.Context, source, event, anonymousID string, properties map[string]any) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.getReporter(source).EventOccurredAnonymous(ctx, event, anonymousID, properties)
}

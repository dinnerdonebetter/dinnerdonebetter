package segment

import (
	"context"
	"errors"

	segment "gopkg.in/segmentio/analytics-go.v3"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	name = "segment_event_reporter"
)

var (
	// ErrEmptyAPIToken indicates an empty API token was provided.
	ErrEmptyAPIToken = errors.New("empty API token")
)

type (
	// EventReporter is a Segment-backed EventReporter.
	EventReporter struct {
		tracer tracing.Tracer
		logger logging.Logger
		client segment.Client
	}
)

// NewSegmentEventReporter returns a new Segment-backed EventReporter.
func NewSegmentEventReporter(logger logging.Logger, tracerProvider tracing.TracerProvider, apiKey string) (*EventReporter, error) {
	if apiKey == "" {
		return nil, ErrEmptyAPIToken
	}

	c := &EventReporter{
		tracer: tracing.NewTracer(tracerProvider.Tracer(name)),
		logger: logging.EnsureLogger(logger).WithName(name),
		client: segment.New(apiKey),
	}

	return c, nil
}

// Close wraps the internal client's Close method.
func (c *EventReporter) Close() {
	if err := c.client.Close(); err != nil {
		c.logger.Error(err, "closing connection")
	}
}

// AddUser upsert's a user's identity.
func (c *EventReporter) AddUser(ctx context.Context, userID string, properties map[string]any) error {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	t := segment.NewTraits()
	for k, v := range properties {
		t.Set(k, v)
	}

	i := segment.NewIntegrations().EnableAll()

	return c.client.Enqueue(segment.Identify{
		UserId:       userID,
		Traits:       t,
		Integrations: i,
	})
}

// EventOccurred associates events with a user.
func (c *EventReporter) EventOccurred(ctx context.Context, event types.CustomerEventType, userID string, properties map[string]any) error {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	p := segment.NewProperties()
	for k, v := range properties {
		p.Set(k, v)
	}

	i := segment.NewIntegrations().EnableAll()

	return c.client.Enqueue(segment.Track{
		Event:        string(event),
		UserId:       userID,
		Properties:   p,
		Integrations: i,
	})
}

package segment

import (
	"context"
	"errors"

	"github.com/prixfixeco/backend/pkg/types"

	"gopkg.in/segmentio/analytics-go.v3"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

const (
	name = "segment_collector"
)

var (
	// ErrEmptyAPIToken indicates an empty API token was provided.
	ErrEmptyAPIToken = errors.New("empty API token")
)

type (
	// CustomerDataCollector is a Segment-backed customerdata.Collector.
	CustomerDataCollector struct {
		tracer tracing.Tracer
		logger logging.Logger
		client analytics.Client
	}
)

// NewSegmentCustomerDataCollector returns a new Segment-backed CustomerDataCollector.
func NewSegmentCustomerDataCollector(logger logging.Logger, tracerProvider tracing.TracerProvider, apiKey string) (*CustomerDataCollector, error) {
	if apiKey == "" {
		return nil, ErrEmptyAPIToken
	}

	c := &CustomerDataCollector{
		tracer: tracing.NewTracer(tracerProvider.Tracer(name)),
		logger: logging.EnsureLogger(logger).WithName(name),
		client: analytics.New(apiKey),
	}

	return c, nil
}

// Close wraps the internal client's Close method.
func (c *CustomerDataCollector) Close() error {
	return c.client.Close()
}

// AddUser upsert's a user's identity.
func (c *CustomerDataCollector) AddUser(ctx context.Context, userID string, properties map[string]interface{}) error {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	t := analytics.NewTraits()
	for k, v := range properties {
		t.Set(k, v)
	}

	i := analytics.NewIntegrations().EnableAll()

	return c.client.Enqueue(analytics.Identify{
		UserId:       userID,
		Traits:       t,
		Integrations: i,
	})
}

// EventOccurred associates events with a user.
func (c *CustomerDataCollector) EventOccurred(ctx context.Context, event types.CustomerEventType, userID string, properties map[string]interface{}) error {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	p := analytics.NewProperties()
	for k, v := range properties {
		p.Set(k, v)
	}

	i := analytics.NewIntegrations().EnableAll()

	return c.client.Enqueue(analytics.Track{
		Event:        string(event),
		UserId:       userID,
		Properties:   p,
		Integrations: i,
	})
}

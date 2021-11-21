package segment

import (
	"context"

	"gopkg.in/segmentio/analytics-go.v3"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	name = "segment_collector"
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
func NewSegmentCustomerDataCollector(logger logging.Logger, apiKey string) (*CustomerDataCollector, error) {
	c := &CustomerDataCollector{
		tracer: tracing.NewTracer(name),
		logger: logging.EnsureLogger(logger).WithName(name),
		client: analytics.New(apiKey),
	}

	return c, nil
}

// Close wraps the internal client's Close method.
func (c *CustomerDataCollector) Close() error {
	return c.client.Close()
}

// Identify upsert's a user's identity.
func (c *CustomerDataCollector) Identify(ctx context.Context, userID string, properties map[string]interface{}) error {
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

// Track associates events with a user.
func (c *CustomerDataCollector) Track(ctx context.Context, event, userID string, properties map[string]interface{}) error {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	p := analytics.NewProperties()
	for k, v := range properties {
		p.Set(k, v)
	}

	i := analytics.NewIntegrations().EnableAll()

	return c.client.Enqueue(analytics.Track{
		Event:        event,
		UserId:       userID,
		Properties:   p,
		Integrations: i,
	})
}

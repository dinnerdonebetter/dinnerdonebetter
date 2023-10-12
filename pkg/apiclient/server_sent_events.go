package apiclient

import (
	"context"

	"github.com/r3labs/sse"
)

// SubscribeToServerSentEventStream subscribes to a server sent event stream.
func (c *Client) SubscribeToServerSentEventStream(ctx context.Context) *sse.Client {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	client := sse.NewClient(c.buildRawURL(ctx, nil, "events").String())
	client.Connection = c.authedClient

	return client
}

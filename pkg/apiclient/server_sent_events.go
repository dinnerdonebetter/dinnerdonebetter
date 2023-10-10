package apiclient

import (
	"context"

	"github.com/r3labs/sse"
)

// SubscribeToServerSentEventStream subscribes to a server sent event stream.
func (c *Client) SubscribeToServerSentEventStream(ctx context.Context) (*sse.Client, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	self, err := c.GetSelf(ctx)
	if err != nil {
		return nil, err
	}

	client := sse.NewClient(c.buildRawURL(ctx, nil, "events").String())
	client.Connection = c.authedClient

	if err = client.SubscribeWithContext(ctx, self.ID, func(msg *sse.Event) {

	}); err != nil {
		return nil, err
	}

	return client, nil
}

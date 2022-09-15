package httpclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetWebhook retrieves a webhook.
func (c *Client) GetWebhook(ctx context.Context, webhookID string) (*types.Webhook, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if webhookID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)

	req, err := c.requestBuilder.BuildGetWebhookRequest(ctx, webhookID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get webhook request")
	}

	var webhook *types.Webhook
	if err = c.fetchAndUnmarshal(ctx, req, &webhook); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving webhook")
	}

	return webhook, nil
}

// GetWebhooks gets a list of webhooks.
func (c *Client) GetWebhooks(ctx context.Context, filter *types.QueryFilter) (*types.WebhookList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetWebhooksRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building webhooks list request")
	}

	var webhooks *types.WebhookList
	if err = c.fetchAndUnmarshal(ctx, req, &webhooks); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving webhooks")
	}

	return webhooks, nil
}

// CreateWebhook creates a webhook.
func (c *Client) CreateWebhook(ctx context.Context, input *types.WebhookCreationRequestInput) (*types.Webhook, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, input.ID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateWebhookRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create webhook request")
	}

	var webhook *types.Webhook
	if err = c.fetchAndUnmarshal(ctx, req, &webhook); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating webhook")
	}

	return webhook, nil
}

// ArchiveWebhook archives a webhook.
func (c *Client) ArchiveWebhook(ctx context.Context, webhookID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if webhookID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)

	req, err := c.requestBuilder.BuildArchiveWebhookRequest(ctx, webhookID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive webhook request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving webhook")
	}

	return nil
}

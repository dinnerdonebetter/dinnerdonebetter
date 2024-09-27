package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
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

	res, err := c.authedGeneratedClient.GetWebhook(ctx, webhookID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting webhook")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Webhook]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading webhook response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetWebhooks gets a list of webhooks.
func (c *Client) GetWebhooks(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Webhook], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetWebhooksParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetWebhooks(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting webhooks")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.Webhook]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading webhooks response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	webhooks := &types.QueryFilteredResult[types.Webhook]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
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

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateWebhookJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateWebhook(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating webhook")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Webhook]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing webhook response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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

	res, err := c.authedGeneratedClient.ArchiveWebhook(ctx, webhookID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving webhook")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Webhook]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading archiving webhook response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// AddWebhookTriggerEvent adds a webhook trigger event.
func (c *Client) AddWebhookTriggerEvent(ctx context.Context, webhookID string, input *types.WebhookTriggerEventCreationRequestInput) (*types.WebhookTriggerEvent, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if webhookID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateWebhookTriggerEventJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateWebhookTriggerEvent(ctx, webhookID, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating webhook trigger event")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.WebhookTriggerEvent]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading webhook trigger event creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// ArchiveWebhookTriggerEvent archives a webhook trigger event.
func (c *Client) ArchiveWebhookTriggerEvent(ctx context.Context, webhookID, webhookTriggerEventID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if webhookID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)

	if webhookTriggerEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookTriggerEventIDKey, webhookTriggerEventID)

	res, err := c.authedGeneratedClient.ArchiveWebhookTriggerEvent(ctx, webhookID, webhookTriggerEventID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving webhook trigger event")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Webhook]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading archiving webhook trigger event response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

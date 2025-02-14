// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

func (c *Client) CreateWebhookTriggerEvent(
	ctx context.Context,
	webhookID string,
	input *WebhookTriggerEventCreationRequestInput,
	reqMods ...RequestModifier,
) (*WebhookTriggerEvent, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if webhookID == "" {
		return nil, buildInvalidIDError("webhook")
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/webhooks/%s/trigger_events", webhookID))
	req, err := c.buildDataRequest(ctx, http.MethodPost, u, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a WebhookTriggerEvent")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*WebhookTriggerEvent]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading WebhookTriggerEvent creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

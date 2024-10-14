// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) GetWebhook(
	ctx context.Context,
	webhookID string,
) (*types.Webhook, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if webhookID == "" {
		return nil, buildInvalidIDError("webhook")
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/webhooks/%s", webhookID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a Webhook")
	}

	var apiResponse *types.APIResponse[*types.Webhook]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading Webhook response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient




import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
)


func (c *Client) CreateWebhookTriggerEvent(
	ctx context.Context,
webhookID string,
input *types.WebhookTriggerEventCreationRequestInput,
) (*types.WebhookTriggerEvent, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}


	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	if webhookID == "" {
		return nil, buildInvalidIDError("webhook")
	} 
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/webhooks/%s/trigger_events" , webhookID ))
	req, err := c.buildDataRequest(ctx, http.MethodPost, u, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a WebhookTriggerEvent")
	}

	var apiResponse *types.APIResponse[ *types.WebhookTriggerEvent]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading WebhookTriggerEvent creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}


	return apiResponse.Data, nil
}
// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) ArchiveOAuth2Client(
	ctx context.Context,
	oauth2ClientID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if oauth2ClientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientIDKey, oauth2ClientID)
	tracing.AttachToSpan(span, keys.OAuth2ClientIDKey, oauth2ClientID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/oauth2_clients/%s", oauth2ClientID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a OAuth2Client")
	}

	var apiResponse *types.APIResponse[*types.OAuth2Client]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading OAuth2Client creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
)

func (c *Client) GetAuthStatus(
	ctx context.Context,
	reqMods ...RequestModifier,
) (*UserStatusResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	u := c.BuildURL(ctx, nil, "/auth/status")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a UserStatusResponse")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*UserStatusResponse]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading UserStatusResponse response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

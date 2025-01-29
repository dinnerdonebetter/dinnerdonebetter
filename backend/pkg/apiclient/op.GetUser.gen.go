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

func (c *Client) GetUser(
	ctx context.Context,
	userID string,
	reqMods ...RequestModifier,
) (*User, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if userID == "" {
		return nil, buildInvalidIDError("user")
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/users/%s", userID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a User")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*User]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading User response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

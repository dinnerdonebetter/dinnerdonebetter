// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

func (c *Client) GetAccount(
	ctx context.Context,
	accountID string,
	reqMods ...RequestModifier,
) (*Account, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if accountID == "" {
		return nil, buildInvalidIDError("account")
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/accounts/%s", accountID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a Account")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*Account]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading Account response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

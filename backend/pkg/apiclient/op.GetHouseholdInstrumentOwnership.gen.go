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

func (c *Client) GetAccountInstrumentOwnership(
	ctx context.Context,
	accountInstrumentOwnershipID string,
	reqMods ...RequestModifier,
) (*AccountInstrumentOwnership, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if accountInstrumentOwnershipID == "" {
		return nil, buildInvalidIDError("accountInstrumentOwnership")
	}
	logger = logger.WithValue(keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/accounts/instruments/%s", accountInstrumentOwnershipID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a AccountInstrumentOwnership")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*AccountInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading AccountInstrumentOwnership response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

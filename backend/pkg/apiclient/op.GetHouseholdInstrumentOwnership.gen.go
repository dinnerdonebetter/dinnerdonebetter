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

func (c *Client) GetHouseholdInstrumentOwnership(
	ctx context.Context,
	householdInstrumentOwnershipID string,
	reqMods ...RequestModifier,
) (*HouseholdInstrumentOwnership, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return nil, buildInvalidIDError("householdInstrumentOwnership")
	}
	logger = logger.WithValue(keys.InstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.InstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/households/instruments/%s", householdInstrumentOwnershipID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a HouseholdInstrumentOwnership")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*HouseholdInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading HouseholdInstrumentOwnership response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

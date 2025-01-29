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

func (c *Client) TransferHouseholdOwnership(
	ctx context.Context,
	householdID string,
	input *HouseholdOwnershipTransferInput,
	reqMods ...RequestModifier,
) (*Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if householdID == "" {
		return nil, buildInvalidIDError("household")
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/households/%s/transfer", householdID))
	req, err := c.buildDataRequest(ctx, http.MethodPost, u, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a Household")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading Household creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

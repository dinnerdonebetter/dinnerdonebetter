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

func (c *Client) CreateHouseholdInvitation(
	ctx context.Context,
	householdID string,
	input *HouseholdInvitationCreationRequestInput,
	reqMods ...RequestModifier,
) (*HouseholdInvitation, error) {
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

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/households/%s/invite", householdID))
	req, err := c.buildDataRequest(ctx, http.MethodPost, u, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a HouseholdInvitation")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading HouseholdInvitation creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

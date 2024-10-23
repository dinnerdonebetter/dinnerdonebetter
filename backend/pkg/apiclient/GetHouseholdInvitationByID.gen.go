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

func (c *Client) GetHouseholdInvitationByID(
	ctx context.Context,
	householdID string,
	householdInvitationID string,
	reqMods ...RequestModifier,
) (*types.HouseholdInvitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdID == "" {
		return nil, buildInvalidIDError("household")
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if householdInvitationID == "" {
		return nil, buildInvalidIDError("householdInvitation")
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/households/%s/invitations/%s", householdID, householdInvitationID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a HouseholdInvitation")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading HouseholdInvitation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

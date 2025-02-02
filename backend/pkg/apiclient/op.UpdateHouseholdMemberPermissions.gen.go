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

func (c *Client) UpdateHouseholdMemberPermissions(
	ctx context.Context,
	householdID string,
	userID string,
	input *ModifyUserPermissionsInput,
	reqMods ...RequestModifier,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/households/%s/members/%s/permissions", householdID, userID))
	req, err := c.buildDataRequest(ctx, http.MethodPatch, u, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a UserPermissionsResponse")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*UserPermissionsResponse]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading UserPermissionsResponse creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

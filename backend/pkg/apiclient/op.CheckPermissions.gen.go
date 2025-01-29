// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
)

func (c *Client) CheckPermissions(
	ctx context.Context,
	input *UserPermissionsRequestInput,
	reqMods ...RequestModifier,
) (*UserPermissionsResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	u := c.BuildURL(ctx, nil, "/api/v1/users/permissions/check")
	req, err := c.buildDataRequest(ctx, http.MethodPost, u, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a UserPermissionsResponse")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*UserPermissionsResponse]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading UserPermissionsResponse creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

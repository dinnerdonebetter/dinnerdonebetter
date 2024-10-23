// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) UpdateUserEmailAddress(
	ctx context.Context,
	input *types.UserEmailAddressUpdateInput,
	reqMods ...RequestModifier,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	u := c.BuildURL(ctx, nil, "/api/v1/users/email_address")
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a User")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[*types.User]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading User creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

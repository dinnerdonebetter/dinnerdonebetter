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

func (c *Client) AcceptAccountInvitation(
	ctx context.Context,
	accountInvitationID string,
	input *AccountInvitationUpdateRequestInput,
	reqMods ...RequestModifier,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if accountInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountInvitationIDKey, accountInvitationID)
	tracing.AttachToSpan(span, keys.AccountInvitationIDKey, accountInvitationID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/account_invitations/%s/accept", accountInvitationID))
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a AccountInvitation")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*AccountInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading AccountInvitation creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

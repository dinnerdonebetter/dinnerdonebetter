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

func (c *Client) GetAccountInvitationByID(
	ctx context.Context,
	accountID string,
	accountInvitationID string,
	reqMods ...RequestModifier,
) (*AccountInvitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if accountID == "" {
		return nil, buildInvalidIDError("account")
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	if accountInvitationID == "" {
		return nil, buildInvalidIDError("accountInvitation")
	}
	logger = logger.WithValue(keys.AccountInvitationIDKey, accountInvitationID)
	tracing.AttachToSpan(span, keys.AccountInvitationIDKey, accountInvitationID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/accounts/%s/invitations/%s", accountID, accountInvitationID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a AccountInvitation")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*AccountInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading AccountInvitation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

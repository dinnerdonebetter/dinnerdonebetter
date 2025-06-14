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

func (c *Client) ArchiveAccountInstrumentOwnership(
	ctx context.Context,
	accountInstrumentOwnershipID string,
	reqMods ...RequestModifier,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if accountInstrumentOwnershipID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/accounts/instruments/%s", accountInstrumentOwnershipID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a AccountInstrumentOwnership")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*AccountInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading AccountInstrumentOwnership creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

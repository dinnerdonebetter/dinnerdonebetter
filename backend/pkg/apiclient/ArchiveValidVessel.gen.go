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

func (c *Client) ArchiveValidVessel(
	ctx context.Context,
	validVesselID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_vessels/%s", validVesselID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a ValidVessel")
	}

	var apiResponse *types.APIResponse[*types.ValidVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading ValidVessel creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

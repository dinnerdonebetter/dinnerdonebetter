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

func (c *Client) GetValidPreparationVessel(
	ctx context.Context,
	validPreparationVesselID string,
	reqMods ...RequestModifier,
) (*ValidPreparationVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVesselID == "" {
		return nil, buildInvalidIDError("validPreparationVessel")
	}
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_preparation_vessels/%s", validPreparationVesselID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a ValidPreparationVessel")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading ValidPreparationVessel response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

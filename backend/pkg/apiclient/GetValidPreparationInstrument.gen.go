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

func (c *Client) GetValidPreparationInstrument(
	ctx context.Context,
	validPreparationVesselID string,
) (*types.ValidPreparationInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVesselID == "" {
		return nil, buildInvalidIDError("validPreparationVessel")
	}
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_preparation_instruments/%s", validPreparationVesselID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a ValidPreparationInstrument")
	}

	var apiResponse *types.APIResponse[*types.ValidPreparationInstrument]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading ValidPreparationInstrument response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

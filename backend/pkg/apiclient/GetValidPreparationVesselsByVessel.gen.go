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

func (c *Client) GetValidPreparationVesselsByVessel(
	ctx context.Context,
	ValidVesselID string,
	filter *types.QueryFilter,
	reqMods ...RequestModifier,
) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if ValidVesselID == "" {
		return nil, buildInvalidIDError("ValidVessel")
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, ValidVesselID)

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, fmt.Sprintf("/api/v1/valid_preparation_vessels/by_vessel/%s", ValidVesselID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of ValidPreparationVessel")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[[]*types.ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of ValidPreparationVessel")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

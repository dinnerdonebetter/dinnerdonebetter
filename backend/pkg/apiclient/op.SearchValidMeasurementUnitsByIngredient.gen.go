// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
)

func (c *Client) SearchValidMeasurementUnitsByIngredient(
	ctx context.Context,
	q string,
	validIngredientID string,
	filter *QueryFilter,
	reqMods ...RequestModifier,
) (*QueryFilteredResult[ValidMeasurementUnit], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = DefaultQueryFilter()
	}
	// tracing.AttachQueryFilterToSpan(span, filter)

	if validIngredientID == "" {
		return nil, buildInvalidIDError("validIngredient")
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	values := filter.ToValues()
	values.Set(textsearch.QueryKeySearch, q)

	u := c.BuildURL(ctx, values, fmt.Sprintf("/api/v1/valid_measurement_units/by_ingredient/%s", validIngredientID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of ValidMeasurementUnit")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[[]*ValidMeasurementUnit]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of ValidMeasurementUnit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &QueryFilteredResult[ValidMeasurementUnit]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

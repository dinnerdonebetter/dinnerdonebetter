// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text"
)

func (c *Client) SearchForValidMeasurementUnits(
	ctx context.Context,
	q string,
	filter *QueryFilter,
	reqMods ...RequestModifier,
) (*QueryFilteredResult[ValidMeasurementUnit], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	values := filter.ToValues()
	values.Set(textsearch.QueryKeySearch, q)

	u := c.BuildURL(ctx, values, "/api/v1/valid_measurement_units/search")
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

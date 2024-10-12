// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient




import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
)


func (c *Client) GetSearchValidMeasurementUnitsByIngredient(
	ctx context.Context,
	q string,
	validIngredientID string,
	filter *types.QueryFilter,
) (*types.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if q == "" {
		return nil, buildInvalidIDError("q")
	} 
	logger = logger.WithValue(keys.SearchQueryKey, q)
	tracing.AttachToSpan(span, keys.SearchQueryKey, q)

	if validIngredientID == "" {
		return nil, buildInvalidIDError("validIngredient")
	} 
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

 

	u := c.BuildURL(ctx, filter.ToValues(), fmt.Sprintf("/api/v1/valid_measurement_units/by_ingredient/%s" , q , validIngredientID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of ValidMeasurementUnit")
	}
	
	var apiResponse *types.APIResponse[ []*types.ValidMeasurementUnit]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of ValidMeasurementUnit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.ValidMeasurementUnit]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}
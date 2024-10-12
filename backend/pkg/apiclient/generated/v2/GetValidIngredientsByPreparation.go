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


func (c *Client) GetValidIngredientsByPreparation(
	ctx context.Context,
	q string,
	validPreparationID string,
	filter *types.QueryFilter,
) (*types.QueryFilteredResult[types.ValidIngredient], error) {
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

	if validPreparationID == "" {
		return nil, buildInvalidIDError("validPreparation")
	} 
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

 

	u := c.BuildURL(ctx, filter.ToValues(), fmt.Sprintf("/api/v1/valid_ingredients/by_preparation/%s" , q , validPreparationID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of ValidIngredient")
	}
	
	var apiResponse *types.APIResponse[ []*types.ValidIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of ValidIngredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.ValidIngredient]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}
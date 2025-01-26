// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) GetValidIngredientStateIngredients(
	ctx context.Context,
	filter *filtering.QueryFilter,
	reqMods ...RequestModifier,
) (*filtering.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, "/api/v1/valid_ingredient_state_ingredients")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of ValidIngredientStateIngredient")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[[]*types.ValidIngredientStateIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of ValidIngredientStateIngredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &filtering.QueryFilteredResult[types.ValidIngredientStateIngredient]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

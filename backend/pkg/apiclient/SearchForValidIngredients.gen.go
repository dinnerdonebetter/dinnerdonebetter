// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) SearchForValidIngredients(
	ctx context.Context,
	q string,
	filter *filtering.QueryFilter,
	reqMods ...RequestModifier,
) (*filtering.QueryFilteredResult[types.ValidIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	values := filter.ToValues()
	values.Set(textsearch.QueryKeySearch, q)

	u := c.BuildURL(ctx, values, "/api/v1/valid_ingredients/search")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of ValidIngredient")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[[]*types.ValidIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of ValidIngredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &filtering.QueryFilteredResult[types.ValidIngredient]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

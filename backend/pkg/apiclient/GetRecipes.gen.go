// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) GetRecipes(
	ctx context.Context,
	filter *types.QueryFilter,
) (*types.QueryFilteredResult[types.Recipe], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, "/api/v1/recipes")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of Recipe")
	}

	var apiResponse *types.APIResponse[[]*types.Recipe]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of Recipe")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.Recipe]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}
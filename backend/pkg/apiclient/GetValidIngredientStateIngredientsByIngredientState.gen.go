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

func (c *Client) GetValidIngredientStateIngredientsByIngredientState(
	ctx context.Context,
	validIngredientStateID string,
	filter *types.QueryFilter,
	reqMods ...RequestModifier,
) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validIngredientStateID == "" {
		return nil, buildInvalidIDError("validIngredientState")
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, fmt.Sprintf("/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/%s", validIngredientStateID))
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

	result := &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

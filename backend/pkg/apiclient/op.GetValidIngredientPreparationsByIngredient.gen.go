// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

func (c *Client) GetValidIngredientPreparationsByIngredient(
	ctx context.Context,
	validIngredientID string,
	filter *QueryFilter,
	reqMods ...RequestModifier,
) (*QueryFilteredResult[ValidIngredientPreparation], error) {
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

	u := c.BuildURL(ctx, values, fmt.Sprintf("/api/v1/valid_ingredient_preparations/by_ingredient/%s", validIngredientID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of ValidIngredientPreparation")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[[]*ValidIngredientPreparation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of ValidIngredientPreparation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &QueryFilteredResult[ValidIngredientPreparation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

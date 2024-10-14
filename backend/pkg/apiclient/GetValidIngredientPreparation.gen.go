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

func (c *Client) GetValidIngredientPreparation(
	ctx context.Context,
	validIngredientPreparationID string,
) (*types.ValidIngredientPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientPreparationID == "" {
		return nil, buildInvalidIDError("validIngredientPreparation")
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_ingredient_preparations/%s", validIngredientPreparationID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a ValidIngredientPreparation")
	}

	var apiResponse *types.APIResponse[*types.ValidIngredientPreparation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading ValidIngredientPreparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

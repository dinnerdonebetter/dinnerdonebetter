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

func (c *Client) GetValidIngredientState(
	ctx context.Context,
	validIngredientStateID string,
	reqMods ...RequestModifier,
) (*types.ValidIngredientState, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateID == "" {
		return nil, buildInvalidIDError("validIngredientState")
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_ingredient_states/%s", validIngredientStateID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a ValidIngredientState")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[*types.ValidIngredientState]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading ValidIngredientState response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

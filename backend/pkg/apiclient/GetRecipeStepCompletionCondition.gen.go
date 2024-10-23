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

func (c *Client) GetRecipeStepCompletionCondition(
	ctx context.Context,
	recipeID string,
	recipeStepID string,
	recipeStepCompletionConditionID string,
	reqMods ...RequestModifier,
) (*types.RecipeStepCompletionCondition, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipeStep")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepCompletionConditionID == "" {
		return nil, buildInvalidIDError("recipeStepCompletionCondition")
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/recipes/%s/steps/%s/completion_conditions/%s", recipeID, recipeStepID, recipeStepCompletionConditionID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a RecipeStepCompletionCondition")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[*types.RecipeStepCompletionCondition]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading RecipeStepCompletionCondition response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

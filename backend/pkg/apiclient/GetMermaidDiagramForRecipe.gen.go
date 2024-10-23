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

func (c *Client) GetMermaidDiagramForRecipe(
	ctx context.Context,
	recipeID string,
	reqMods ...RequestModifier,
) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return "", buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/recipes/%s/mermaid", recipeID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "building request to fetch a string")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[string]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "loading string response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return "", err
	}

	return apiResponse.Data, nil
}

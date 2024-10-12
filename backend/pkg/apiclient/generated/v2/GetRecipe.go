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


func (c *Client) GetRecipe(
	ctx context.Context,
recipeID string,
) ( *types.Recipe, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	} 
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/recipes/%s" , recipeID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a Recipe")
	}

	var apiResponse *types.APIResponse[  *types.Recipe]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading Recipe response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}


	return apiResponse.Data, nil
}
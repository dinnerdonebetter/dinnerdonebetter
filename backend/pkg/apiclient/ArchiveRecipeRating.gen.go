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

func (c *Client) ArchiveRecipeRating(
	ctx context.Context,
	recipeID string,
	recipeRatingID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeRatingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/recipes/%s/ratings/%s", recipeID, recipeRatingID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a RecipeRating")
	}

	var apiResponse *types.APIResponse[*types.RecipeRating]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading RecipeRating creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

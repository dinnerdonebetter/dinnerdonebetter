package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// UploadRecipeMedia uploads a piece of media for a recipe.
// TODO: write unit test for this.
func (c *Client) UploadRecipeMedia(ctx context.Context, files map[string][]byte, recipeID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return buildInvalidIDError("recipe")
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if files == nil {
		return ErrNilInputProvided
	}

	uri := c.BuildURL(ctx, nil, "recipes", recipeID, "images")

	req, err := c.buildMultipleRecipeMediaUploadRequest(ctx, uri, files)
	if err != nil {
		return observability.PrepareError(err, span, "media upload")
	}

	var apiResponse *types.APIResponse[[]*types.RecipeMedia]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "uploading media")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

// UploadMediaForRecipeStep uploads a piece of media for a recipe step.
// TODO: write unit test for this.
func (c *Client) UploadMediaForRecipeStep(ctx context.Context, files map[string][]byte, recipeID, recipeStepID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if files == nil {
		return ErrNilInputProvided
	}

	uri := c.BuildURL(ctx, nil, "recipes", recipeID, "steps", recipeStepID, "images")

	logger := c.logger.WithValue(keys.RecipeStepIDKey, recipeID).WithValue(keys.RecipeStepIDKey, recipeStepID)
	logger.WithValue("uri", uri).Info("Uploading recipe step media")

	req, err := c.buildMultipleRecipeMediaUploadRequest(ctx, uri, files)
	if err != nil {
		return observability.PrepareError(err, span, "recipe step media upload")
	}

	var apiResponse *APIResponse[[]*RecipeMedia]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "uploading recipe step media")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

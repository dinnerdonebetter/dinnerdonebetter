package requests

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

const (
	recipeStepsBasePath = "steps"
)

// BuildMultipleRecipeMediaUploadRequestForRecipeStep builds an HTTP request that sets a user's avatar to the provided content.
func (b *Builder) BuildMultipleRecipeMediaUploadRequestForRecipeStep(ctx context.Context, files map[string][]byte, recipeID, recipeStepID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	body, formDataContentType, err := b.prepareUploads(ctx, files)
	if err != nil {
		return nil, observability.PrepareError(err, span, "preparing upload request")
	}

	uri := b.BuildURL(ctx, nil, recipesBasePath, recipeID, recipeStepsBasePath, recipeStepID, "images")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building media upload request")
	}

	req.Header.Set("Content-Type", formDataContentType)

	return req, nil
}

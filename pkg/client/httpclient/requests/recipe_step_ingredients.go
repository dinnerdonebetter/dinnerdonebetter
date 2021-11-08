package requests

import (
	"context"
	"net/http"

	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	recipeStepIngredientsBasePath = "recipe_step_ingredients"
)

// BuildGetRecipeStepIngredientRequest builds an HTTP request for fetching a recipe step ingredient.
func (b *Builder) BuildGetRecipeStepIngredientRequest(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepIngredientsBasePath,
		recipeStepIngredientID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetRecipeStepIngredientsRequest builds an HTTP request for fetching a list of recipe step ingredients.
func (b *Builder) BuildGetRecipeStepIngredientsRequest(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateRecipeStepIngredientRequest builds an HTTP request for creating a recipe step ingredient.
func (b *Builder) BuildCreateRecipeStepIngredientRequest(ctx context.Context, recipeID string, input *types.RecipeStepIngredientCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		input.BelongsToRecipeStep,
		recipeStepIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateRecipeStepIngredientRequest builds an HTTP request for updating a recipe step ingredient.
func (b *Builder) BuildUpdateRecipeStepIngredientRequest(ctx context.Context, recipeID string, recipeStepIngredient *types.RecipeStepIngredient) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepIngredient == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredient.ID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredient.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepIngredient.BelongsToRecipeStep,
		recipeStepIngredientsBasePath,
		recipeStepIngredient.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, recipeStepIngredient)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveRecipeStepIngredientRequest builds an HTTP request for archiving a recipe step ingredient.
func (b *Builder) BuildArchiveRecipeStepIngredientRequest(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepIngredientsBasePath,
		recipeStepIngredientID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

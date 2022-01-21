package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	recipesBasePath = "recipes"
)

// BuildGetRecipeRequest builds an HTTP request for fetching a recipe.
func (b *Builder) BuildGetRecipeRequest(ctx context.Context, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetRecipesRequest builds an HTTP request for fetching a list of recipes.
func (b *Builder) BuildGetRecipesRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		recipesBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildSearchForRecipesRequest builds an HTTP request for fetching a list of recipes.
func (b *Builder) BuildSearchForRecipesRequest(ctx context.Context, query string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger).WithValue(keys.SearchQueryKey, query)
	tracing.AttachSearchQueryToSpan(span, query)

	queryParams := filter.ToValues()
	queryParams.Set(types.SearchQueryKey, query)

	uri := b.BuildURL(
		ctx,
		queryParams,
		recipesBasePath,
		"search",
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateRecipeRequest builds an HTTP request for creating a recipe.
func (b *Builder) BuildCreateRecipeRequest(ctx context.Context, input *types.RecipeCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

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
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateRecipeRequest builds an HTTP request for updating a recipe.
func (b *Builder) BuildUpdateRecipeRequest(ctx context.Context, recipe *types.Recipe) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if recipe == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.RecipeIDKey, recipe.ID)
	tracing.AttachRecipeIDToSpan(span, recipe.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipe.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, recipe)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveRecipeRequest builds an HTTP request for archiving a recipe.
func (b *Builder) BuildArchiveRecipeRequest(ctx context.Context, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

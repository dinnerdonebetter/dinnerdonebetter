package requests

import (
	"context"
	"net/http"
	"strconv"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	recipeStepProductsBasePath = "recipe_step_products"
)

// BuildRecipeStepProductExistsRequest builds an HTTP request for checking the existence of a recipe step product.
func (b *Builder) BuildRecipeStepProductExistsRequest(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		id(recipeID),
		recipeStepsBasePath,
		id(recipeStepID),
		recipeStepProductsBasePath,
		id(recipeStepProductID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetRecipeStepProductRequest builds an HTTP request for fetching a recipe step product.
func (b *Builder) BuildGetRecipeStepProductRequest(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		id(recipeID),
		recipeStepsBasePath,
		id(recipeStepID),
		recipeStepProductsBasePath,
		id(recipeStepProductID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetRecipeStepProductsRequest builds an HTTP request for fetching a list of recipe step products.
func (b *Builder) BuildGetRecipeStepProductsRequest(ctx context.Context, recipeID, recipeStepID uint64, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		recipesBasePath,
		id(recipeID),
		recipeStepsBasePath,
		id(recipeStepID),
		recipeStepProductsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateRecipeStepProductRequest builds an HTTP request for creating a recipe step product.
func (b *Builder) BuildCreateRecipeStepProductRequest(ctx context.Context, recipeID uint64, input *types.RecipeStepProductCreationInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if recipeID == 0 {
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
		id(recipeID),
		recipeStepsBasePath,
		id(input.BelongsToRecipeStep),
		recipeStepProductsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateRecipeStepProductRequest builds an HTTP request for updating a recipe step product.
func (b *Builder) BuildUpdateRecipeStepProductRequest(ctx context.Context, recipeID uint64, recipeStepProduct *types.RecipeStepProduct) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepProduct == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProduct.ID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProduct.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		id(recipeID),
		recipeStepsBasePath,
		id(recipeStepProduct.BelongsToRecipeStep),
		recipeStepProductsBasePath,
		strconv.FormatUint(recipeStepProduct.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, recipeStepProduct)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveRecipeStepProductRequest builds an HTTP request for archiving a recipe step product.
func (b *Builder) BuildArchiveRecipeStepProductRequest(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		id(recipeID),
		recipeStepsBasePath,
		id(recipeStepID),
		recipeStepProductsBasePath,
		id(recipeStepProductID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetAuditLogForRecipeStepProductRequest builds an HTTP request for fetching a list of audit log entries pertaining to a recipe step product.
func (b *Builder) BuildGetAuditLogForRecipeStepProductRequest(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		id(recipeID),
		recipeStepsBasePath,
		id(recipeStepID),
		recipeStepProductsBasePath,
		id(recipeStepProductID),
		"audit",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

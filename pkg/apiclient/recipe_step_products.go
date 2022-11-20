package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

// GetRecipeStepProduct gets a recipe step product.
func (c *Client) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	req, err := c.requestBuilder.BuildGetRecipeStepProductRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe step product request")
	}

	var recipeStepProduct *types.RecipeStepProduct
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepProduct); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step product")
	}

	return recipeStepProduct, nil
}

// GetRecipeStepProducts retrieves a list of recipe step products.
func (c *Client) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepProduct], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	req, err := c.requestBuilder.BuildGetRecipeStepProductsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe step products list request")
	}

	var recipeStepProducts *types.QueryFilteredResult[types.RecipeStepProduct]
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepProducts); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step products")
	}

	return recipeStepProducts, nil
}

// CreateRecipeStepProduct creates a recipe step product.
func (c *Client) CreateRecipeStepProduct(ctx context.Context, recipeID string, input *types.RecipeStepProductCreationRequestInput) (*types.RecipeStepProduct, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepProductRequest(ctx, recipeID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe step product request")
	}

	var recipeStepProduct *types.RecipeStepProduct
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepProduct); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step product")
	}

	return recipeStepProduct, nil
}

// UpdateRecipeStepProduct updates a recipe step product.
func (c *Client) UpdateRecipeStepProduct(ctx context.Context, recipeID string, recipeStepProduct *types.RecipeStepProduct) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepProduct == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProduct.ID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProduct.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepProductRequest(ctx, recipeID, recipeStepProduct)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe step product request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepProduct); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product %s", recipeStepProduct.ID)
	}

	return nil
}

// ArchiveRecipeStepProduct archives a recipe step product.
func (c *Client) ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepProductRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe step product request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step product %s", recipeStepProductID)
	}

	return nil
}

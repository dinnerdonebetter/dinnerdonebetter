package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// RecipeStepProductExists retrieves whether a recipe step product exists.
func (c *Client) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	req, err := c.requestBuilder.BuildRecipeStepProductExistsRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building recipe step product existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for recipe step product #%d", recipeStepProductID)
	}

	return exists, nil
}

// GetRecipeStepProduct gets a recipe step product.
func (c *Client) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*types.RecipeStepProduct, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

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

	req, err := c.requestBuilder.BuildGetRecipeStepProductRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get recipe step product request")
	}

	var recipeStepProduct *types.RecipeStepProduct
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepProduct); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving recipe step product")
	}

	return recipeStepProduct, nil
}

// GetRecipeStepProducts retrieves a list of recipe step products.
func (c *Client) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *types.QueryFilter) (*types.RecipeStepProductList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

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

	req, err := c.requestBuilder.BuildGetRecipeStepProductsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building recipe step products list request")
	}

	var recipeStepProducts *types.RecipeStepProductList
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepProducts); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving recipe step products")
	}

	return recipeStepProducts, nil
}

// CreateRecipeStepProduct creates a recipe step product.
func (c *Client) CreateRecipeStepProduct(ctx context.Context, recipeID uint64, input *types.RecipeStepProductCreationInput) (*types.RecipeStepProduct, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

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

	req, err := c.requestBuilder.BuildCreateRecipeStepProductRequest(ctx, recipeID, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create recipe step product request")
	}

	var recipeStepProduct *types.RecipeStepProduct
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepProduct); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating recipe step product")
	}

	return recipeStepProduct, nil
}

// UpdateRecipeStepProduct updates a recipe step product.
func (c *Client) UpdateRecipeStepProduct(ctx context.Context, recipeID uint64, recipeStepProduct *types.RecipeStepProduct) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
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
		return observability.PrepareError(err, logger, span, "building update recipe step product request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepProduct); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step product #%d", recipeStepProduct.ID)
	}

	return nil
}

// ArchiveRecipeStepProduct archives a recipe step product.
func (c *Client) ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepProductRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive recipe step product request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving recipe step product #%d", recipeStepProductID)
	}

	return nil
}

// GetAuditLogForRecipeStepProduct retrieves a list of audit log entries pertaining to a recipe step product.
func (c *Client) GetAuditLogForRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

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

	req, err := c.requestBuilder.BuildGetAuditLogForRecipeStepProductRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for recipe step product request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

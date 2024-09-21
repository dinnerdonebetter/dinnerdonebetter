package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepProductID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	res, err := c.authedGeneratedClient.GetRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe step product")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepProduct]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step product")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRecipeStepProducts retrieves a list of recipe step products.
func (c *Client) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepProduct], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	params := &generated.GetRecipeStepProductsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetRecipeStepProducts(ctx, recipeID, recipeStepID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "recipe step products list")
	}

	var apiResponse *types.APIResponse[[]*types.RecipeStepProduct]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step products")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.RecipeStepProduct]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateRecipeStepProduct creates a recipe step product.
func (c *Client) CreateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductCreationRequestInput) (*types.RecipeStepProduct, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeStepID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateRecipeStepProductJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateRecipeStepProduct(ctx, recipeID, recipeStepID, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create recipe step product")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepProduct]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step product")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepProduct == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProduct.ID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProduct.ID)

	body := generated.UpdateRecipeStepProductJSONRequestBody{}
	c.copyType(&body, recipeStepProduct)

	res, err := c.authedGeneratedClient.UpdateRecipeStepProduct(ctx, recipeID, recipeStepProduct.BelongsToRecipeStep, recipeStepProduct.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update recipe step product")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepProduct]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product %s", recipeStepProduct.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepProductID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)

	res, err := c.authedGeneratedClient.ArchiveRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive recipe step product")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepProduct]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step product %s", recipeStepProductID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

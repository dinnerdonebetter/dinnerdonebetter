package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipeStepProductsBasePath = "recipe_step_products"
)

// BuildRecipeStepProductExistsRequest builds an HTTP request for checking the existence of a recipe step product.
func (c *V1Client) BuildRecipeStepProductExistsRequest(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeStepProductExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepProductsBasePath,
		strconv.FormatUint(recipeStepProductID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeStepProductExists retrieves whether or not a recipe step product exists.
func (c *V1Client) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepProductExists")
	defer span.End()

	req, err := c.BuildRecipeStepProductExistsRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeStepProductRequest builds an HTTP request for fetching a recipe step product.
func (c *V1Client) BuildGetRecipeStepProductRequest(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepProductRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepProductsBasePath,
		strconv.FormatUint(recipeStepProductID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepProduct retrieves a recipe step product.
func (c *V1Client) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (recipeStepProduct *models.RecipeStepProduct, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepProduct")
	defer span.End()

	req, err := c.BuildGetRecipeStepProductRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepProduct); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepProduct, nil
}

// BuildGetRecipeStepProductsRequest builds an HTTP request for fetching recipe step products.
func (c *V1Client) BuildGetRecipeStepProductsRequest(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepProductsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepProductsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepProducts retrieves a list of recipe step products.
func (c *V1Client) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (recipeStepProducts *models.RecipeStepProductList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepProducts")
	defer span.End()

	req, err := c.BuildGetRecipeStepProductsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepProducts); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepProducts, nil
}

// BuildCreateRecipeStepProductRequest builds an HTTP request for creating a recipe step product.
func (c *V1Client) BuildCreateRecipeStepProductRequest(ctx context.Context, recipeID uint64, input *models.RecipeStepProductCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeStepProductRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(input.BelongsToRecipeStep, 10),
		recipeStepProductsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeStepProduct creates a recipe step product.
func (c *V1Client) CreateRecipeStepProduct(ctx context.Context, recipeID uint64, input *models.RecipeStepProductCreationInput) (recipeStepProduct *models.RecipeStepProduct, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepProduct")
	defer span.End()

	req, err := c.BuildCreateRecipeStepProductRequest(ctx, recipeID, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepProduct)
	return recipeStepProduct, err
}

// BuildUpdateRecipeStepProductRequest builds an HTTP request for updating a recipe step product.
func (c *V1Client) BuildUpdateRecipeStepProductRequest(ctx context.Context, recipeID uint64, recipeStepProduct *models.RecipeStepProduct) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeStepProductRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepProduct.BelongsToRecipeStep, 10),
		recipeStepProductsBasePath,
		strconv.FormatUint(recipeStepProduct.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeStepProduct)
}

// UpdateRecipeStepProduct updates a recipe step product.
func (c *V1Client) UpdateRecipeStepProduct(ctx context.Context, recipeID uint64, recipeStepProduct *models.RecipeStepProduct) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepProduct")
	defer span.End()

	req, err := c.BuildUpdateRecipeStepProductRequest(ctx, recipeID, recipeStepProduct)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeStepProduct)
}

// BuildArchiveRecipeStepProductRequest builds an HTTP request for updating a recipe step product.
func (c *V1Client) BuildArchiveRecipeStepProductRequest(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeStepProductRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepProductsBasePath,
		strconv.FormatUint(recipeStepProductID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepProduct archives a recipe step product.
func (c *V1Client) ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepProduct")
	defer span.End()

	req, err := c.BuildArchiveRecipeStepProductRequest(ctx, recipeID, recipeStepID, recipeStepProductID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

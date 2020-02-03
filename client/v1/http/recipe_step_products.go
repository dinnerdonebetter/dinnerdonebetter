package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipeStepProductsBasePath = "recipe_step_products"
)

// BuildGetRecipeStepProductRequest builds an HTTP request for fetching a recipe step product
func (c *V1Client) BuildGetRecipeStepProductRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepProductsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStepProduct retrieves a recipe step product
func (c *V1Client) GetRecipeStepProduct(ctx context.Context, id uint64) (recipeStepProduct *models.RecipeStepProduct, err error) {
	req, err := c.BuildGetRecipeStepProductRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepProduct); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepProduct, nil
}

// BuildGetRecipeStepProductsRequest builds an HTTP request for fetching recipe step products
func (c *V1Client) BuildGetRecipeStepProductsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), recipeStepProductsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStepProducts retrieves a list of recipe step products
func (c *V1Client) GetRecipeStepProducts(ctx context.Context, filter *models.QueryFilter) (recipeStepProducts *models.RecipeStepProductList, err error) {
	req, err := c.BuildGetRecipeStepProductsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepProducts); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepProducts, nil
}

// BuildCreateRecipeStepProductRequest builds an HTTP request for creating a recipe step product
func (c *V1Client) BuildCreateRecipeStepProductRequest(ctx context.Context, body *models.RecipeStepProductCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepProductsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateRecipeStepProduct creates a recipe step product
func (c *V1Client) CreateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProductCreationInput) (recipeStepProduct *models.RecipeStepProduct, err error) {
	req, err := c.BuildCreateRecipeStepProductRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepProduct)
	return recipeStepProduct, err
}

// BuildUpdateRecipeStepProductRequest builds an HTTP request for updating a recipe step product
func (c *V1Client) BuildUpdateRecipeStepProductRequest(ctx context.Context, updated *models.RecipeStepProduct) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepProductsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateRecipeStepProduct updates a recipe step product
func (c *V1Client) UpdateRecipeStepProduct(ctx context.Context, updated *models.RecipeStepProduct) error {
	req, err := c.BuildUpdateRecipeStepProductRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveRecipeStepProductRequest builds an HTTP request for updating a recipe step product
func (c *V1Client) BuildArchiveRecipeStepProductRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepProductsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepProduct archives a recipe step product
func (c *V1Client) ArchiveRecipeStepProduct(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveRecipeStepProductRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

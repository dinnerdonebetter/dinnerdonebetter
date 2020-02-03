package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	ingredientsBasePath = "ingredients"
)

// BuildGetIngredientRequest builds an HTTP request for fetching an ingredient
func (c *V1Client) BuildGetIngredientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, ingredientsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetIngredient retrieves an ingredient
func (c *V1Client) GetIngredient(ctx context.Context, id uint64) (ingredient *models.Ingredient, err error) {
	req, err := c.BuildGetIngredientRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &ingredient); retrieveErr != nil {
		return nil, retrieveErr
	}

	return ingredient, nil
}

// BuildGetIngredientsRequest builds an HTTP request for fetching ingredients
func (c *V1Client) BuildGetIngredientsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), ingredientsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetIngredients retrieves a list of ingredients
func (c *V1Client) GetIngredients(ctx context.Context, filter *models.QueryFilter) (ingredients *models.IngredientList, err error) {
	req, err := c.BuildGetIngredientsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &ingredients); retrieveErr != nil {
		return nil, retrieveErr
	}

	return ingredients, nil
}

// BuildCreateIngredientRequest builds an HTTP request for creating an ingredient
func (c *V1Client) BuildCreateIngredientRequest(ctx context.Context, body *models.IngredientCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, ingredientsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateIngredient creates an ingredient
func (c *V1Client) CreateIngredient(ctx context.Context, input *models.IngredientCreationInput) (ingredient *models.Ingredient, err error) {
	req, err := c.BuildCreateIngredientRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &ingredient)
	return ingredient, err
}

// BuildUpdateIngredientRequest builds an HTTP request for updating an ingredient
func (c *V1Client) BuildUpdateIngredientRequest(ctx context.Context, updated *models.Ingredient) (*http.Request, error) {
	uri := c.BuildURL(nil, ingredientsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateIngredient updates an ingredient
func (c *V1Client) UpdateIngredient(ctx context.Context, updated *models.Ingredient) error {
	req, err := c.BuildUpdateIngredientRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveIngredientRequest builds an HTTP request for updating an ingredient
func (c *V1Client) BuildArchiveIngredientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, ingredientsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveIngredient archives an ingredient
func (c *V1Client) ArchiveIngredient(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveIngredientRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipeStepIngredientsBasePath = "recipe_step_ingredients"
)

// BuildGetRecipeStepIngredientRequest builds an HTTP request for fetching a recipe step ingredient
func (c *V1Client) BuildGetRecipeStepIngredientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepIngredientsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStepIngredient retrieves a recipe step ingredient
func (c *V1Client) GetRecipeStepIngredient(ctx context.Context, id uint64) (recipeStepIngredient *models.RecipeStepIngredient, err error) {
	req, err := c.BuildGetRecipeStepIngredientRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepIngredient); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepIngredient, nil
}

// BuildGetRecipeStepIngredientsRequest builds an HTTP request for fetching recipe step ingredients
func (c *V1Client) BuildGetRecipeStepIngredientsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), recipeStepIngredientsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStepIngredients retrieves a list of recipe step ingredients
func (c *V1Client) GetRecipeStepIngredients(ctx context.Context, filter *models.QueryFilter) (recipeStepIngredients *models.RecipeStepIngredientList, err error) {
	req, err := c.BuildGetRecipeStepIngredientsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepIngredients); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepIngredients, nil
}

// BuildCreateRecipeStepIngredientRequest builds an HTTP request for creating a recipe step ingredient
func (c *V1Client) BuildCreateRecipeStepIngredientRequest(ctx context.Context, body *models.RecipeStepIngredientCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepIngredientsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateRecipeStepIngredient creates a recipe step ingredient
func (c *V1Client) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (recipeStepIngredient *models.RecipeStepIngredient, err error) {
	req, err := c.BuildCreateRecipeStepIngredientRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepIngredient)
	return recipeStepIngredient, err
}

// BuildUpdateRecipeStepIngredientRequest builds an HTTP request for updating a recipe step ingredient
func (c *V1Client) BuildUpdateRecipeStepIngredientRequest(ctx context.Context, updated *models.RecipeStepIngredient) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepIngredientsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateRecipeStepIngredient updates a recipe step ingredient
func (c *V1Client) UpdateRecipeStepIngredient(ctx context.Context, updated *models.RecipeStepIngredient) error {
	req, err := c.BuildUpdateRecipeStepIngredientRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveRecipeStepIngredientRequest builds an HTTP request for updating a recipe step ingredient
func (c *V1Client) BuildArchiveRecipeStepIngredientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepIngredientsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient
func (c *V1Client) ArchiveRecipeStepIngredient(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveRecipeStepIngredientRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipesBasePath = "recipes"
)

// BuildGetRecipeRequest builds an HTTP request for fetching a recipe
func (c *V1Client) BuildGetRecipeRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipesBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipe retrieves a recipe
func (c *V1Client) GetRecipe(ctx context.Context, id uint64) (recipe *models.Recipe, err error) {
	req, err := c.BuildGetRecipeRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipe); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipe, nil
}

// BuildGetRecipesRequest builds an HTTP request for fetching recipes
func (c *V1Client) BuildGetRecipesRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), recipesBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipes retrieves a list of recipes
func (c *V1Client) GetRecipes(ctx context.Context, filter *models.QueryFilter) (recipes *models.RecipeList, err error) {
	req, err := c.BuildGetRecipesRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipes); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipes, nil
}

// BuildCreateRecipeRequest builds an HTTP request for creating a recipe
func (c *V1Client) BuildCreateRecipeRequest(ctx context.Context, body *models.RecipeCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, recipesBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateRecipe creates a recipe
func (c *V1Client) CreateRecipe(ctx context.Context, input *models.RecipeCreationInput) (recipe *models.Recipe, err error) {
	req, err := c.BuildCreateRecipeRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipe)
	return recipe, err
}

// BuildUpdateRecipeRequest builds an HTTP request for updating a recipe
func (c *V1Client) BuildUpdateRecipeRequest(ctx context.Context, updated *models.Recipe) (*http.Request, error) {
	uri := c.BuildURL(nil, recipesBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateRecipe updates a recipe
func (c *V1Client) UpdateRecipe(ctx context.Context, updated *models.Recipe) error {
	req, err := c.BuildUpdateRecipeRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveRecipeRequest builds an HTTP request for updating a recipe
func (c *V1Client) BuildArchiveRecipeRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipesBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveRecipe archives a recipe
func (c *V1Client) ArchiveRecipe(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveRecipeRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

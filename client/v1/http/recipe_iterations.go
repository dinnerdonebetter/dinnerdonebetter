package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipeIterationsBasePath = "recipe_iterations"
)

// BuildGetRecipeIterationRequest builds an HTTP request for fetching a recipe iteration
func (c *V1Client) BuildGetRecipeIterationRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeIterationsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeIteration retrieves a recipe iteration
func (c *V1Client) GetRecipeIteration(ctx context.Context, id uint64) (recipeIteration *models.RecipeIteration, err error) {
	req, err := c.BuildGetRecipeIterationRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeIteration); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeIteration, nil
}

// BuildGetRecipeIterationsRequest builds an HTTP request for fetching recipe iterations
func (c *V1Client) BuildGetRecipeIterationsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), recipeIterationsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeIterations retrieves a list of recipe iterations
func (c *V1Client) GetRecipeIterations(ctx context.Context, filter *models.QueryFilter) (recipeIterations *models.RecipeIterationList, err error) {
	req, err := c.BuildGetRecipeIterationsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeIterations); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeIterations, nil
}

// BuildCreateRecipeIterationRequest builds an HTTP request for creating a recipe iteration
func (c *V1Client) BuildCreateRecipeIterationRequest(ctx context.Context, body *models.RecipeIterationCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeIterationsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateRecipeIteration creates a recipe iteration
func (c *V1Client) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (recipeIteration *models.RecipeIteration, err error) {
	req, err := c.BuildCreateRecipeIterationRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeIteration)
	return recipeIteration, err
}

// BuildUpdateRecipeIterationRequest builds an HTTP request for updating a recipe iteration
func (c *V1Client) BuildUpdateRecipeIterationRequest(ctx context.Context, updated *models.RecipeIteration) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeIterationsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateRecipeIteration updates a recipe iteration
func (c *V1Client) UpdateRecipeIteration(ctx context.Context, updated *models.RecipeIteration) error {
	req, err := c.BuildUpdateRecipeIterationRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveRecipeIterationRequest builds an HTTP request for updating a recipe iteration
func (c *V1Client) BuildArchiveRecipeIterationRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeIterationsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveRecipeIteration archives a recipe iteration
func (c *V1Client) ArchiveRecipeIteration(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveRecipeIterationRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

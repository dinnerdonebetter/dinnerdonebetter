package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipeStepsBasePath = "recipe_steps"
)

// BuildGetRecipeStepRequest builds an HTTP request for fetching a recipe step
func (c *V1Client) BuildGetRecipeStepRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStep retrieves a recipe step
func (c *V1Client) GetRecipeStep(ctx context.Context, id uint64) (recipeStep *models.RecipeStep, err error) {
	req, err := c.BuildGetRecipeStepRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStep); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStep, nil
}

// BuildGetRecipeStepsRequest builds an HTTP request for fetching recipe steps
func (c *V1Client) BuildGetRecipeStepsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), recipeStepsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeSteps retrieves a list of recipe steps
func (c *V1Client) GetRecipeSteps(ctx context.Context, filter *models.QueryFilter) (recipeSteps *models.RecipeStepList, err error) {
	req, err := c.BuildGetRecipeStepsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeSteps); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeSteps, nil
}

// BuildCreateRecipeStepRequest builds an HTTP request for creating a recipe step
func (c *V1Client) BuildCreateRecipeStepRequest(ctx context.Context, body *models.RecipeStepCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateRecipeStep creates a recipe step
func (c *V1Client) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (recipeStep *models.RecipeStep, err error) {
	req, err := c.BuildCreateRecipeStepRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStep)
	return recipeStep, err
}

// BuildUpdateRecipeStepRequest builds an HTTP request for updating a recipe step
func (c *V1Client) BuildUpdateRecipeStepRequest(ctx context.Context, updated *models.RecipeStep) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateRecipeStep updates a recipe step
func (c *V1Client) UpdateRecipeStep(ctx context.Context, updated *models.RecipeStep) error {
	req, err := c.BuildUpdateRecipeStepRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveRecipeStepRequest builds an HTTP request for updating a recipe step
func (c *V1Client) BuildArchiveRecipeStepRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveRecipeStep archives a recipe step
func (c *V1Client) ArchiveRecipeStep(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveRecipeStepRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

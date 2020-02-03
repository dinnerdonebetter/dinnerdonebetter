package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipeStepEventsBasePath = "recipe_step_events"
)

// BuildGetRecipeStepEventRequest builds an HTTP request for fetching a recipe step event
func (c *V1Client) BuildGetRecipeStepEventRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepEventsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStepEvent retrieves a recipe step event
func (c *V1Client) GetRecipeStepEvent(ctx context.Context, id uint64) (recipeStepEvent *models.RecipeStepEvent, err error) {
	req, err := c.BuildGetRecipeStepEventRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepEvent); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepEvent, nil
}

// BuildGetRecipeStepEventsRequest builds an HTTP request for fetching recipe step events
func (c *V1Client) BuildGetRecipeStepEventsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), recipeStepEventsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStepEvents retrieves a list of recipe step events
func (c *V1Client) GetRecipeStepEvents(ctx context.Context, filter *models.QueryFilter) (recipeStepEvents *models.RecipeStepEventList, err error) {
	req, err := c.BuildGetRecipeStepEventsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepEvents); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepEvents, nil
}

// BuildCreateRecipeStepEventRequest builds an HTTP request for creating a recipe step event
func (c *V1Client) BuildCreateRecipeStepEventRequest(ctx context.Context, body *models.RecipeStepEventCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepEventsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateRecipeStepEvent creates a recipe step event
func (c *V1Client) CreateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEventCreationInput) (recipeStepEvent *models.RecipeStepEvent, err error) {
	req, err := c.BuildCreateRecipeStepEventRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepEvent)
	return recipeStepEvent, err
}

// BuildUpdateRecipeStepEventRequest builds an HTTP request for updating a recipe step event
func (c *V1Client) BuildUpdateRecipeStepEventRequest(ctx context.Context, updated *models.RecipeStepEvent) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepEventsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateRecipeStepEvent updates a recipe step event
func (c *V1Client) UpdateRecipeStepEvent(ctx context.Context, updated *models.RecipeStepEvent) error {
	req, err := c.BuildUpdateRecipeStepEventRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveRecipeStepEventRequest builds an HTTP request for updating a recipe step event
func (c *V1Client) BuildArchiveRecipeStepEventRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepEventsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepEvent archives a recipe step event
func (c *V1Client) ArchiveRecipeStepEvent(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveRecipeStepEventRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

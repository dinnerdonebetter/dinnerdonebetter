package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipeStepInstrumentsBasePath = "recipe_step_instruments"
)

// BuildGetRecipeStepInstrumentRequest builds an HTTP request for fetching a recipe step instrument
func (c *V1Client) BuildGetRecipeStepInstrumentRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepInstrumentsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStepInstrument retrieves a recipe step instrument
func (c *V1Client) GetRecipeStepInstrument(ctx context.Context, id uint64) (recipeStepInstrument *models.RecipeStepInstrument, err error) {
	req, err := c.BuildGetRecipeStepInstrumentRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepInstrument); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepInstrument, nil
}

// BuildGetRecipeStepInstrumentsRequest builds an HTTP request for fetching recipe step instruments
func (c *V1Client) BuildGetRecipeStepInstrumentsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), recipeStepInstrumentsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRecipeStepInstruments retrieves a list of recipe step instruments
func (c *V1Client) GetRecipeStepInstruments(ctx context.Context, filter *models.QueryFilter) (recipeStepInstruments *models.RecipeStepInstrumentList, err error) {
	req, err := c.BuildGetRecipeStepInstrumentsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepInstruments); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepInstruments, nil
}

// BuildCreateRecipeStepInstrumentRequest builds an HTTP request for creating a recipe step instrument
func (c *V1Client) BuildCreateRecipeStepInstrumentRequest(ctx context.Context, body *models.RecipeStepInstrumentCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepInstrumentsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateRecipeStepInstrument creates a recipe step instrument
func (c *V1Client) CreateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrumentCreationInput) (recipeStepInstrument *models.RecipeStepInstrument, err error) {
	req, err := c.BuildCreateRecipeStepInstrumentRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepInstrument)
	return recipeStepInstrument, err
}

// BuildUpdateRecipeStepInstrumentRequest builds an HTTP request for updating a recipe step instrument
func (c *V1Client) BuildUpdateRecipeStepInstrumentRequest(ctx context.Context, updated *models.RecipeStepInstrument) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepInstrumentsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateRecipeStepInstrument updates a recipe step instrument
func (c *V1Client) UpdateRecipeStepInstrument(ctx context.Context, updated *models.RecipeStepInstrument) error {
	req, err := c.BuildUpdateRecipeStepInstrumentRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveRecipeStepInstrumentRequest builds an HTTP request for updating a recipe step instrument
func (c *V1Client) BuildArchiveRecipeStepInstrumentRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, recipeStepInstrumentsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepInstrument archives a recipe step instrument
func (c *V1Client) ArchiveRecipeStepInstrument(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveRecipeStepInstrumentRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

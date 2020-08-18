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
	recipeStepInstrumentsBasePath = "recipe_step_instruments"
)

// BuildRecipeStepInstrumentExistsRequest builds an HTTP request for checking the existence of a recipe step instrument.
func (c *V1Client) BuildRecipeStepInstrumentExistsRequest(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeStepInstrumentExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepInstrumentsBasePath,
		strconv.FormatUint(recipeStepInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeStepInstrumentExists retrieves whether or not a recipe step instrument exists.
func (c *V1Client) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepInstrumentExists")
	defer span.End()

	req, err := c.BuildRecipeStepInstrumentExistsRequest(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeStepInstrumentRequest builds an HTTP request for fetching a recipe step instrument.
func (c *V1Client) BuildGetRecipeStepInstrumentRequest(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepInstrumentsBasePath,
		strconv.FormatUint(recipeStepInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepInstrument retrieves a recipe step instrument.
func (c *V1Client) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (recipeStepInstrument *models.RecipeStepInstrument, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepInstrument")
	defer span.End()

	req, err := c.BuildGetRecipeStepInstrumentRequest(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepInstrument); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepInstrument, nil
}

// BuildGetRecipeStepInstrumentsRequest builds an HTTP request for fetching recipe step instruments.
func (c *V1Client) BuildGetRecipeStepInstrumentsRequest(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepInstrumentsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepInstruments retrieves a list of recipe step instruments.
func (c *V1Client) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (recipeStepInstruments *models.RecipeStepInstrumentList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepInstruments")
	defer span.End()

	req, err := c.BuildGetRecipeStepInstrumentsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepInstruments); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepInstruments, nil
}

// BuildCreateRecipeStepInstrumentRequest builds an HTTP request for creating a recipe step instrument.
func (c *V1Client) BuildCreateRecipeStepInstrumentRequest(ctx context.Context, recipeID uint64, input *models.RecipeStepInstrumentCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeStepInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(input.BelongsToRecipeStep, 10),
		recipeStepInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeStepInstrument creates a recipe step instrument.
func (c *V1Client) CreateRecipeStepInstrument(ctx context.Context, recipeID uint64, input *models.RecipeStepInstrumentCreationInput) (recipeStepInstrument *models.RecipeStepInstrument, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepInstrument")
	defer span.End()

	req, err := c.BuildCreateRecipeStepInstrumentRequest(ctx, recipeID, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepInstrument)
	return recipeStepInstrument, err
}

// BuildUpdateRecipeStepInstrumentRequest builds an HTTP request for updating a recipe step instrument.
func (c *V1Client) BuildUpdateRecipeStepInstrumentRequest(ctx context.Context, recipeID uint64, recipeStepInstrument *models.RecipeStepInstrument) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeStepInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepInstrument.BelongsToRecipeStep, 10),
		recipeStepInstrumentsBasePath,
		strconv.FormatUint(recipeStepInstrument.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeStepInstrument)
}

// UpdateRecipeStepInstrument updates a recipe step instrument.
func (c *V1Client) UpdateRecipeStepInstrument(ctx context.Context, recipeID uint64, recipeStepInstrument *models.RecipeStepInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepInstrument")
	defer span.End()

	req, err := c.BuildUpdateRecipeStepInstrumentRequest(ctx, recipeID, recipeStepInstrument)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeStepInstrument)
}

// BuildArchiveRecipeStepInstrumentRequest builds an HTTP request for updating a recipe step instrument.
func (c *V1Client) BuildArchiveRecipeStepInstrumentRequest(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeStepInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepInstrumentsBasePath,
		strconv.FormatUint(recipeStepInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepInstrument archives a recipe step instrument.
func (c *V1Client) ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepInstrument")
	defer span.End()

	req, err := c.BuildArchiveRecipeStepInstrumentRequest(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

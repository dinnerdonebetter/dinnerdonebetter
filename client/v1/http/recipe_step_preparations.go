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
	recipeStepPreparationsBasePath = "recipe_step_preparations"
)

// BuildRecipeStepPreparationExistsRequest builds an HTTP request for checking the existence of a recipe step preparation.
func (c *V1Client) BuildRecipeStepPreparationExistsRequest(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeStepPreparationExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepPreparationsBasePath,
		strconv.FormatUint(recipeStepPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeStepPreparationExists retrieves whether or not a recipe step preparation exists.
func (c *V1Client) RecipeStepPreparationExists(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepPreparationExists")
	defer span.End()

	req, err := c.BuildRecipeStepPreparationExistsRequest(ctx, recipeID, recipeStepID, recipeStepPreparationID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeStepPreparationRequest builds an HTTP request for fetching a recipe step preparation.
func (c *V1Client) BuildGetRecipeStepPreparationRequest(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepPreparationsBasePath,
		strconv.FormatUint(recipeStepPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepPreparation retrieves a recipe step preparation.
func (c *V1Client) GetRecipeStepPreparation(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (recipeStepPreparation *models.RecipeStepPreparation, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepPreparation")
	defer span.End()

	req, err := c.BuildGetRecipeStepPreparationRequest(ctx, recipeID, recipeStepID, recipeStepPreparationID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepPreparation); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepPreparation, nil
}

// BuildGetRecipeStepPreparationsRequest builds an HTTP request for fetching recipe step preparations.
func (c *V1Client) BuildGetRecipeStepPreparationsRequest(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepPreparationsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepPreparations retrieves a list of recipe step preparations.
func (c *V1Client) GetRecipeStepPreparations(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (recipeStepPreparations *models.RecipeStepPreparationList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepPreparations")
	defer span.End()

	req, err := c.BuildGetRecipeStepPreparationsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepPreparations); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepPreparations, nil
}

// BuildCreateRecipeStepPreparationRequest builds an HTTP request for creating a recipe step preparation.
func (c *V1Client) BuildCreateRecipeStepPreparationRequest(ctx context.Context, recipeID uint64, input *models.RecipeStepPreparationCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeStepPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(input.BelongsToRecipeStep, 10),
		recipeStepPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeStepPreparation creates a recipe step preparation.
func (c *V1Client) CreateRecipeStepPreparation(ctx context.Context, recipeID uint64, input *models.RecipeStepPreparationCreationInput) (recipeStepPreparation *models.RecipeStepPreparation, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepPreparation")
	defer span.End()

	req, err := c.BuildCreateRecipeStepPreparationRequest(ctx, recipeID, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepPreparation)
	return recipeStepPreparation, err
}

// BuildUpdateRecipeStepPreparationRequest builds an HTTP request for updating a recipe step preparation.
func (c *V1Client) BuildUpdateRecipeStepPreparationRequest(ctx context.Context, recipeID uint64, recipeStepPreparation *models.RecipeStepPreparation) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeStepPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepPreparation.BelongsToRecipeStep, 10),
		recipeStepPreparationsBasePath,
		strconv.FormatUint(recipeStepPreparation.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeStepPreparation)
}

// UpdateRecipeStepPreparation updates a recipe step preparation.
func (c *V1Client) UpdateRecipeStepPreparation(ctx context.Context, recipeID uint64, recipeStepPreparation *models.RecipeStepPreparation) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepPreparation")
	defer span.End()

	req, err := c.BuildUpdateRecipeStepPreparationRequest(ctx, recipeID, recipeStepPreparation)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeStepPreparation)
}

// BuildArchiveRecipeStepPreparationRequest builds an HTTP request for updating a recipe step preparation.
func (c *V1Client) BuildArchiveRecipeStepPreparationRequest(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeStepPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepPreparationsBasePath,
		strconv.FormatUint(recipeStepPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepPreparation archives a recipe step preparation.
func (c *V1Client) ArchiveRecipeStepPreparation(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepPreparation")
	defer span.End()

	req, err := c.BuildArchiveRecipeStepPreparationRequest(ctx, recipeID, recipeStepID, recipeStepPreparationID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

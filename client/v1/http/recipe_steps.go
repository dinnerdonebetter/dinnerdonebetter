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
	recipeStepsBasePath = "recipe_steps"
)

// BuildRecipeStepExistsRequest builds an HTTP request for checking the existence of a recipe step.
func (c *V1Client) BuildRecipeStepExistsRequest(ctx context.Context, recipeID, recipeStepID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeStepExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeStepExists retrieves whether or not a recipe step exists.
func (c *V1Client) RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepExists")
	defer span.End()

	req, err := c.BuildRecipeStepExistsRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeStepRequest builds an HTTP request for fetching a recipe step.
func (c *V1Client) BuildGetRecipeStepRequest(ctx context.Context, recipeID, recipeStepID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStep retrieves a recipe step.
func (c *V1Client) GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (recipeStep *models.RecipeStep, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStep")
	defer span.End()

	req, err := c.BuildGetRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStep); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStep, nil
}

// BuildGetRecipeStepsRequest builds an HTTP request for fetching recipe steps.
func (c *V1Client) BuildGetRecipeStepsRequest(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeSteps retrieves a list of recipe steps.
func (c *V1Client) GetRecipeSteps(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (recipeSteps *models.RecipeStepList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeSteps")
	defer span.End()

	req, err := c.BuildGetRecipeStepsRequest(ctx, recipeID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeSteps); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeSteps, nil
}

// BuildCreateRecipeStepRequest builds an HTTP request for creating a recipe step.
func (c *V1Client) BuildCreateRecipeStepRequest(ctx context.Context, input *models.RecipeStepCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeStepRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(input.BelongsToRecipe, 10),
		recipeStepsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeStep creates a recipe step.
func (c *V1Client) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (recipeStep *models.RecipeStep, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStep")
	defer span.End()

	req, err := c.BuildCreateRecipeStepRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStep)
	return recipeStep, err
}

// BuildUpdateRecipeStepRequest builds an HTTP request for updating a recipe step.
func (c *V1Client) BuildUpdateRecipeStepRequest(ctx context.Context, recipeStep *models.RecipeStep) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeStepRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeStep.BelongsToRecipe, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStep.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeStep)
}

// UpdateRecipeStep updates a recipe step.
func (c *V1Client) UpdateRecipeStep(ctx context.Context, recipeStep *models.RecipeStep) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStep")
	defer span.End()

	req, err := c.BuildUpdateRecipeStepRequest(ctx, recipeStep)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeStep)
}

// BuildArchiveRecipeStepRequest builds an HTTP request for updating a recipe step.
func (c *V1Client) BuildArchiveRecipeStepRequest(ctx context.Context, recipeID, recipeStepID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeStepRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeStep archives a recipe step.
func (c *V1Client) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStep")
	defer span.End()

	req, err := c.BuildArchiveRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

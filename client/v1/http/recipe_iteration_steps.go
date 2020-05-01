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
	recipeIterationStepsBasePath = "recipe_iteration_steps"
)

// BuildRecipeIterationStepExistsRequest builds an HTTP request for checking the existence of a recipe iteration step.
func (c *V1Client) BuildRecipeIterationStepExistsRequest(ctx context.Context, recipeID, recipeIterationStepID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeIterationStepExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationStepsBasePath,
		strconv.FormatUint(recipeIterationStepID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeIterationStepExists retrieves whether or not a recipe iteration step exists.
func (c *V1Client) RecipeIterationStepExists(ctx context.Context, recipeID, recipeIterationStepID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeIterationStepExists")
	defer span.End()

	req, err := c.BuildRecipeIterationStepExistsRequest(ctx, recipeID, recipeIterationStepID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeIterationStepRequest builds an HTTP request for fetching a recipe iteration step.
func (c *V1Client) BuildGetRecipeIterationStepRequest(ctx context.Context, recipeID, recipeIterationStepID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeIterationStepRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationStepsBasePath,
		strconv.FormatUint(recipeIterationStepID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeIterationStep retrieves a recipe iteration step.
func (c *V1Client) GetRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) (recipeIterationStep *models.RecipeIterationStep, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeIterationStep")
	defer span.End()

	req, err := c.BuildGetRecipeIterationStepRequest(ctx, recipeID, recipeIterationStepID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeIterationStep); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeIterationStep, nil
}

// BuildGetRecipeIterationStepsRequest builds an HTTP request for fetching recipe iteration steps.
func (c *V1Client) BuildGetRecipeIterationStepsRequest(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeIterationStepsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationStepsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeIterationSteps retrieves a list of recipe iteration steps.
func (c *V1Client) GetRecipeIterationSteps(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (recipeIterationSteps *models.RecipeIterationStepList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeIterationSteps")
	defer span.End()

	req, err := c.BuildGetRecipeIterationStepsRequest(ctx, recipeID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeIterationSteps); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeIterationSteps, nil
}

// BuildCreateRecipeIterationStepRequest builds an HTTP request for creating a recipe iteration step.
func (c *V1Client) BuildCreateRecipeIterationStepRequest(ctx context.Context, input *models.RecipeIterationStepCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeIterationStepRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(input.BelongsToRecipe, 10),
		recipeIterationStepsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeIterationStep creates a recipe iteration step.
func (c *V1Client) CreateRecipeIterationStep(ctx context.Context, input *models.RecipeIterationStepCreationInput) (recipeIterationStep *models.RecipeIterationStep, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeIterationStep")
	defer span.End()

	req, err := c.BuildCreateRecipeIterationStepRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeIterationStep)
	return recipeIterationStep, err
}

// BuildUpdateRecipeIterationStepRequest builds an HTTP request for updating a recipe iteration step.
func (c *V1Client) BuildUpdateRecipeIterationStepRequest(ctx context.Context, recipeIterationStep *models.RecipeIterationStep) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeIterationStepRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeIterationStep.BelongsToRecipe, 10),
		recipeIterationStepsBasePath,
		strconv.FormatUint(recipeIterationStep.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeIterationStep)
}

// UpdateRecipeIterationStep updates a recipe iteration step.
func (c *V1Client) UpdateRecipeIterationStep(ctx context.Context, recipeIterationStep *models.RecipeIterationStep) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeIterationStep")
	defer span.End()

	req, err := c.BuildUpdateRecipeIterationStepRequest(ctx, recipeIterationStep)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeIterationStep)
}

// BuildArchiveRecipeIterationStepRequest builds an HTTP request for updating a recipe iteration step.
func (c *V1Client) BuildArchiveRecipeIterationStepRequest(ctx context.Context, recipeID, recipeIterationStepID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeIterationStepRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationStepsBasePath,
		strconv.FormatUint(recipeIterationStepID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeIterationStep archives a recipe iteration step.
func (c *V1Client) ArchiveRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeIterationStep")
	defer span.End()

	req, err := c.BuildArchiveRecipeIterationStepRequest(ctx, recipeID, recipeIterationStepID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

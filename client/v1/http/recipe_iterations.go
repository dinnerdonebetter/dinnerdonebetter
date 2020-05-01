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
	recipeIterationsBasePath = "recipe_iterations"
)

// BuildRecipeIterationExistsRequest builds an HTTP request for checking the existence of a recipe iteration.
func (c *V1Client) BuildRecipeIterationExistsRequest(ctx context.Context, recipeID, recipeIterationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeIterationExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(recipeIterationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeIterationExists retrieves whether or not a recipe iteration exists.
func (c *V1Client) RecipeIterationExists(ctx context.Context, recipeID, recipeIterationID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeIterationExists")
	defer span.End()

	req, err := c.BuildRecipeIterationExistsRequest(ctx, recipeID, recipeIterationID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeIterationRequest builds an HTTP request for fetching a recipe iteration.
func (c *V1Client) BuildGetRecipeIterationRequest(ctx context.Context, recipeID, recipeIterationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeIterationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(recipeIterationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeIteration retrieves a recipe iteration.
func (c *V1Client) GetRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) (recipeIteration *models.RecipeIteration, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeIteration")
	defer span.End()

	req, err := c.BuildGetRecipeIterationRequest(ctx, recipeID, recipeIterationID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeIteration); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeIteration, nil
}

// BuildGetRecipeIterationsRequest builds an HTTP request for fetching recipe iterations.
func (c *V1Client) BuildGetRecipeIterationsRequest(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeIterationsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeIterations retrieves a list of recipe iterations.
func (c *V1Client) GetRecipeIterations(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (recipeIterations *models.RecipeIterationList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeIterations")
	defer span.End()

	req, err := c.BuildGetRecipeIterationsRequest(ctx, recipeID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeIterations); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeIterations, nil
}

// BuildCreateRecipeIterationRequest builds an HTTP request for creating a recipe iteration.
func (c *V1Client) BuildCreateRecipeIterationRequest(ctx context.Context, input *models.RecipeIterationCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeIterationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(input.BelongsToRecipe, 10),
		recipeIterationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeIteration creates a recipe iteration.
func (c *V1Client) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (recipeIteration *models.RecipeIteration, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeIteration")
	defer span.End()

	req, err := c.BuildCreateRecipeIterationRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeIteration)
	return recipeIteration, err
}

// BuildUpdateRecipeIterationRequest builds an HTTP request for updating a recipe iteration.
func (c *V1Client) BuildUpdateRecipeIterationRequest(ctx context.Context, recipeIteration *models.RecipeIteration) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeIterationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeIteration.BelongsToRecipe, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(recipeIteration.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeIteration)
}

// UpdateRecipeIteration updates a recipe iteration.
func (c *V1Client) UpdateRecipeIteration(ctx context.Context, recipeIteration *models.RecipeIteration) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeIteration")
	defer span.End()

	req, err := c.BuildUpdateRecipeIterationRequest(ctx, recipeIteration)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeIteration)
}

// BuildArchiveRecipeIterationRequest builds an HTTP request for updating a recipe iteration.
func (c *V1Client) BuildArchiveRecipeIterationRequest(ctx context.Context, recipeID, recipeIterationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeIterationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(recipeIterationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeIteration archives a recipe iteration.
func (c *V1Client) ArchiveRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeIteration")
	defer span.End()

	req, err := c.BuildArchiveRecipeIterationRequest(ctx, recipeID, recipeIterationID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

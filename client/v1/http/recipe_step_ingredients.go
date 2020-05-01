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
	recipeStepIngredientsBasePath = "recipe_step_ingredients"
)

// BuildRecipeStepIngredientExistsRequest builds an HTTP request for checking the existence of a recipe step ingredient.
func (c *V1Client) BuildRecipeStepIngredientExistsRequest(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeStepIngredientExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepIngredientsBasePath,
		strconv.FormatUint(recipeStepIngredientID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeStepIngredientExists retrieves whether or not a recipe step ingredient exists.
func (c *V1Client) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepIngredientExists")
	defer span.End()

	req, err := c.BuildRecipeStepIngredientExistsRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeStepIngredientRequest builds an HTTP request for fetching a recipe step ingredient.
func (c *V1Client) BuildGetRecipeStepIngredientRequest(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepIngredientRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepIngredientsBasePath,
		strconv.FormatUint(recipeStepIngredientID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepIngredient retrieves a recipe step ingredient.
func (c *V1Client) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (recipeStepIngredient *models.RecipeStepIngredient, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepIngredient")
	defer span.End()

	req, err := c.BuildGetRecipeStepIngredientRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepIngredient); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepIngredient, nil
}

// BuildGetRecipeStepIngredientsRequest builds an HTTP request for fetching recipe step ingredients.
func (c *V1Client) BuildGetRecipeStepIngredientsRequest(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepIngredientsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepIngredients retrieves a list of recipe step ingredients.
func (c *V1Client) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (recipeStepIngredients *models.RecipeStepIngredientList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepIngredients")
	defer span.End()

	req, err := c.BuildGetRecipeStepIngredientsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepIngredients); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepIngredients, nil
}

// BuildCreateRecipeStepIngredientRequest builds an HTTP request for creating a recipe step ingredient.
func (c *V1Client) BuildCreateRecipeStepIngredientRequest(ctx context.Context, recipeID uint64, input *models.RecipeStepIngredientCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeStepIngredientRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(input.BelongsToRecipeStep, 10),
		recipeStepIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeStepIngredient creates a recipe step ingredient.
func (c *V1Client) CreateRecipeStepIngredient(ctx context.Context, recipeID uint64, input *models.RecipeStepIngredientCreationInput) (recipeStepIngredient *models.RecipeStepIngredient, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepIngredient")
	defer span.End()

	req, err := c.BuildCreateRecipeStepIngredientRequest(ctx, recipeID, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepIngredient)
	return recipeStepIngredient, err
}

// BuildUpdateRecipeStepIngredientRequest builds an HTTP request for updating a recipe step ingredient.
func (c *V1Client) BuildUpdateRecipeStepIngredientRequest(ctx context.Context, recipeID uint64, recipeStepIngredient *models.RecipeStepIngredient) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeStepIngredientRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepIngredient.BelongsToRecipeStep, 10),
		recipeStepIngredientsBasePath,
		strconv.FormatUint(recipeStepIngredient.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeStepIngredient)
}

// UpdateRecipeStepIngredient updates a recipe step ingredient.
func (c *V1Client) UpdateRecipeStepIngredient(ctx context.Context, recipeID uint64, recipeStepIngredient *models.RecipeStepIngredient) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepIngredient")
	defer span.End()

	req, err := c.BuildUpdateRecipeStepIngredientRequest(ctx, recipeID, recipeStepIngredient)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeStepIngredient)
}

// BuildArchiveRecipeStepIngredientRequest builds an HTTP request for updating a recipe step ingredient.
func (c *V1Client) BuildArchiveRecipeStepIngredientRequest(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeStepIngredientRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepIngredientsBasePath,
		strconv.FormatUint(recipeStepIngredientID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient.
func (c *V1Client) ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepIngredient")
	defer span.End()

	req, err := c.BuildArchiveRecipeStepIngredientRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

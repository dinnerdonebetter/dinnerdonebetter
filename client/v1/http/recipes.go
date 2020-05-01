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
	recipesBasePath = "recipes"
)

// BuildRecipeExistsRequest builds an HTTP request for checking the existence of a recipe.
func (c *V1Client) BuildRecipeExistsRequest(ctx context.Context, recipeID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeExists retrieves whether or not a recipe exists.
func (c *V1Client) RecipeExists(ctx context.Context, recipeID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeExists")
	defer span.End()

	req, err := c.BuildRecipeExistsRequest(ctx, recipeID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeRequest builds an HTTP request for fetching a recipe.
func (c *V1Client) BuildGetRecipeRequest(ctx context.Context, recipeID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipe retrieves a recipe.
func (c *V1Client) GetRecipe(ctx context.Context, recipeID uint64) (recipe *models.Recipe, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipe")
	defer span.End()

	req, err := c.BuildGetRecipeRequest(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipe); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipe, nil
}

// BuildGetRecipesRequest builds an HTTP request for fetching recipes.
func (c *V1Client) BuildGetRecipesRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipesRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipes retrieves a list of recipes.
func (c *V1Client) GetRecipes(ctx context.Context, filter *models.QueryFilter) (recipes *models.RecipeList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipes")
	defer span.End()

	req, err := c.BuildGetRecipesRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipes); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipes, nil
}

// BuildCreateRecipeRequest builds an HTTP request for creating a recipe.
func (c *V1Client) BuildCreateRecipeRequest(ctx context.Context, input *models.RecipeCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipe creates a recipe.
func (c *V1Client) CreateRecipe(ctx context.Context, input *models.RecipeCreationInput) (recipe *models.Recipe, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipe")
	defer span.End()

	req, err := c.BuildCreateRecipeRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipe)
	return recipe, err
}

// BuildUpdateRecipeRequest builds an HTTP request for updating a recipe.
func (c *V1Client) BuildUpdateRecipeRequest(ctx context.Context, recipe *models.Recipe) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipe.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipe)
}

// UpdateRecipe updates a recipe.
func (c *V1Client) UpdateRecipe(ctx context.Context, recipe *models.Recipe) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipe")
	defer span.End()

	req, err := c.BuildUpdateRecipeRequest(ctx, recipe)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipe)
}

// BuildArchiveRecipeRequest builds an HTTP request for updating a recipe.
func (c *V1Client) BuildArchiveRecipeRequest(ctx context.Context, recipeID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipe archives a recipe.
func (c *V1Client) ArchiveRecipe(ctx context.Context, recipeID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipe")
	defer span.End()

	req, err := c.BuildArchiveRecipeRequest(ctx, recipeID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

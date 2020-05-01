package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeDataManager = (*Client)(nil)

// RecipeExists fetches whether or not a recipe exists from the database.
func (c *Client) RecipeExists(ctx context.Context, recipeID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
	}).Debug("RecipeExists called")

	return c.querier.RecipeExists(ctx, recipeID)
}

// GetRecipe fetches a recipe from the database.
func (c *Client) GetRecipe(ctx context.Context, recipeID uint64) (*models.Recipe, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipe")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
	}).Debug("GetRecipe called")

	return c.querier.GetRecipe(ctx, recipeID)
}

// GetAllRecipesCount fetches the count of recipes from the database that meet a particular filter.
func (c *Client) GetAllRecipesCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipesCount")
	defer span.End()

	c.logger.Debug("GetAllRecipesCount called")

	return c.querier.GetAllRecipesCount(ctx)
}

// GetRecipes fetches a list of recipes from the database that meet a particular filter.
func (c *Client) GetRecipes(ctx context.Context, filter *models.QueryFilter) (*models.RecipeList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipes")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetRecipes called")

	recipeList, err := c.querier.GetRecipes(ctx, filter)

	return recipeList, err
}

// CreateRecipe creates a recipe in the database.
func (c *Client) CreateRecipe(ctx context.Context, input *models.RecipeCreationInput) (*models.Recipe, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipe")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipe called")

	return c.querier.CreateRecipe(ctx, input)
}

// UpdateRecipe updates a particular recipe. Note that UpdateRecipe expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipe(ctx context.Context, updated *models.Recipe) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipe")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_id", updated.ID).Debug("UpdateRecipe called")

	return c.querier.UpdateRecipe(ctx, updated)
}

// ArchiveRecipe archives a recipe from the database by its ID.
func (c *Client) ArchiveRecipe(ctx context.Context, recipeID, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipe")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
		"user_id":   userID,
	}).Debug("ArchiveRecipe called")

	return c.querier.ArchiveRecipe(ctx, recipeID, userID)
}

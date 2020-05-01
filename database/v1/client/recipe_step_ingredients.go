package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeStepIngredientDataManager = (*Client)(nil)

// RecipeStepIngredientExists fetches whether or not a recipe step ingredient exists from the database.
func (c *Client) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepIngredientExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":                 recipeID,
		"recipe_step_id":            recipeStepID,
		"recipe_step_ingredient_id": recipeStepIngredientID,
	}).Debug("RecipeStepIngredientExists called")

	return c.querier.RecipeStepIngredientExists(ctx, recipeID, recipeStepID, recipeStepIngredientID)
}

// GetRecipeStepIngredient fetches a recipe step ingredient from the database.
func (c *Client) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*models.RecipeStepIngredient, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepIngredient")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":                 recipeID,
		"recipe_step_id":            recipeStepID,
		"recipe_step_ingredient_id": recipeStepIngredientID,
	}).Debug("GetRecipeStepIngredient called")

	return c.querier.GetRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID)
}

// GetAllRecipeStepIngredientsCount fetches the count of recipe step ingredients from the database that meet a particular filter.
func (c *Client) GetAllRecipeStepIngredientsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepIngredientsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepIngredientsCount called")

	return c.querier.GetAllRecipeStepIngredientsCount(ctx)
}

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter.
func (c *Client) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepIngredientList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepIngredients")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":      recipeID,
		"recipe_step_id": recipeStepID,
	}).Debug("GetRecipeStepIngredients called")

	recipeStepIngredientList, err := c.querier.GetRecipeStepIngredients(ctx, recipeID, recipeStepID, filter)

	return recipeStepIngredientList, err
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database.
func (c *Client) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (*models.RecipeStepIngredient, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepIngredient")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepIngredient called")

	return c.querier.CreateRecipeStepIngredient(ctx, input)
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient. Note that UpdateRecipeStepIngredient expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepIngredient(ctx context.Context, updated *models.RecipeStepIngredient) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepIngredient")
	defer span.End()

	tracing.AttachRecipeStepIngredientIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_step_ingredient_id", updated.ID).Debug("UpdateRecipeStepIngredient called")

	return c.querier.UpdateRecipeStepIngredient(ctx, updated)
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient from the database by its ID.
func (c *Client) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepIngredient")
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_ingredient_id": recipeStepIngredientID,
		"recipe_step_id":            recipeStepID,
	}).Debug("ArchiveRecipeStepIngredient called")

	return c.querier.ArchiveRecipeStepIngredient(ctx, recipeStepID, recipeStepIngredientID)
}

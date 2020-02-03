package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.RecipeStepIngredientDataManager = (*Client)(nil)

// attachRecipeStepIngredientIDToSpan provides a consistent way to attach a recipe step ingredient's ID to a span
func attachRecipeStepIngredientIDToSpan(span *trace.Span, recipeStepIngredientID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_step_ingredient_id", strconv.FormatUint(recipeStepIngredientID, 10)))
	}
}

// GetRecipeStepIngredient fetches a recipe step ingredient from the database
func (c *Client) GetRecipeStepIngredient(ctx context.Context, recipeStepIngredientID, userID uint64) (*models.RecipeStepIngredient, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepIngredient")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_ingredient_id": recipeStepIngredientID,
		"user_id":                   userID,
	}).Debug("GetRecipeStepIngredient called")

	return c.querier.GetRecipeStepIngredient(ctx, recipeStepIngredientID, userID)
}

// GetRecipeStepIngredientCount fetches the count of recipe step ingredients from the database that meet a particular filter
func (c *Client) GetRecipeStepIngredientCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepIngredientCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepIngredientCount called")

	return c.querier.GetRecipeStepIngredientCount(ctx, filter, userID)
}

// GetAllRecipeStepIngredientsCount fetches the count of recipe step ingredients from the database that meet a particular filter
func (c *Client) GetAllRecipeStepIngredientsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepIngredientsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepIngredientsCount called")

	return c.querier.GetAllRecipeStepIngredientsCount(ctx)
}

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter
func (c *Client) GetRecipeStepIngredients(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepIngredientList, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepIngredients")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepIngredients called")

	recipeStepIngredientList, err := c.querier.GetRecipeStepIngredients(ctx, filter, userID)

	return recipeStepIngredientList, err
}

// GetAllRecipeStepIngredientsForUser fetches a list of recipe step ingredients from the database that meet a particular filter
func (c *Client) GetAllRecipeStepIngredientsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepIngredient, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepIngredientsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllRecipeStepIngredientsForUser called")

	recipeStepIngredientList, err := c.querier.GetAllRecipeStepIngredientsForUser(ctx, userID)

	return recipeStepIngredientList, err
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database
func (c *Client) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (*models.RecipeStepIngredient, error) {
	ctx, span := trace.StartSpan(ctx, "CreateRecipeStepIngredient")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepIngredient called")

	return c.querier.CreateRecipeStepIngredient(ctx, input)
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient. Note that UpdateRecipeStepIngredient expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredient) error {
	ctx, span := trace.StartSpan(ctx, "UpdateRecipeStepIngredient")
	defer span.End()

	attachRecipeStepIngredientIDToSpan(span, input.ID)
	c.logger.WithValue("recipe_step_ingredient_id", input.ID).Debug("UpdateRecipeStepIngredient called")

	return c.querier.UpdateRecipeStepIngredient(ctx, input)
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient from the database by its ID
func (c *Client) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepIngredientID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveRecipeStepIngredient")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_ingredient_id": recipeStepIngredientID,
		"user_id":                   userID,
	}).Debug("ArchiveRecipeStepIngredient called")

	return c.querier.ArchiveRecipeStepIngredient(ctx, recipeStepIngredientID, userID)
}

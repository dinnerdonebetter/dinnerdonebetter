package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.RecipeDataManager = (*Client)(nil)

// attachRecipeIDToSpan provides a consistent way to attach a recipe's ID to a span
func attachRecipeIDToSpan(span *trace.Span, recipeID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_id", strconv.FormatUint(recipeID, 10)))
	}
}

// GetRecipe fetches a recipe from the database
func (c *Client) GetRecipe(ctx context.Context, recipeID, userID uint64) (*models.Recipe, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipe")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeIDToSpan(span, recipeID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
		"user_id":   userID,
	}).Debug("GetRecipe called")

	return c.querier.GetRecipe(ctx, recipeID, userID)
}

// GetRecipeCount fetches the count of recipes from the database that meet a particular filter
func (c *Client) GetRecipeCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeCount called")

	return c.querier.GetRecipeCount(ctx, filter, userID)
}

// GetAllRecipesCount fetches the count of recipes from the database that meet a particular filter
func (c *Client) GetAllRecipesCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipesCount")
	defer span.End()

	c.logger.Debug("GetAllRecipesCount called")

	return c.querier.GetAllRecipesCount(ctx)
}

// GetRecipes fetches a list of recipes from the database that meet a particular filter
func (c *Client) GetRecipes(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeList, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipes")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipes called")

	recipeList, err := c.querier.GetRecipes(ctx, filter, userID)

	return recipeList, err
}

// GetAllRecipesForUser fetches a list of recipes from the database that meet a particular filter
func (c *Client) GetAllRecipesForUser(ctx context.Context, userID uint64) ([]models.Recipe, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipesForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllRecipesForUser called")

	recipeList, err := c.querier.GetAllRecipesForUser(ctx, userID)

	return recipeList, err
}

// CreateRecipe creates a recipe in the database
func (c *Client) CreateRecipe(ctx context.Context, input *models.RecipeCreationInput) (*models.Recipe, error) {
	ctx, span := trace.StartSpan(ctx, "CreateRecipe")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipe called")

	return c.querier.CreateRecipe(ctx, input)
}

// UpdateRecipe updates a particular recipe. Note that UpdateRecipe expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipe(ctx context.Context, input *models.Recipe) error {
	ctx, span := trace.StartSpan(ctx, "UpdateRecipe")
	defer span.End()

	attachRecipeIDToSpan(span, input.ID)
	c.logger.WithValue("recipe_id", input.ID).Debug("UpdateRecipe called")

	return c.querier.UpdateRecipe(ctx, input)
}

// ArchiveRecipe archives a recipe from the database by its ID
func (c *Client) ArchiveRecipe(ctx context.Context, recipeID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveRecipe")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeIDToSpan(span, recipeID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
		"user_id":   userID,
	}).Debug("ArchiveRecipe called")

	return c.querier.ArchiveRecipe(ctx, recipeID, userID)
}

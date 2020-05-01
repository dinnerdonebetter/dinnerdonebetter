package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeTagDataManager = (*Client)(nil)

// RecipeTagExists fetches whether or not a recipe tag exists from the database.
func (c *Client) RecipeTagExists(ctx context.Context, recipeID, recipeTagID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeTagExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeTagIDToSpan(span, recipeTagID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":     recipeID,
		"recipe_tag_id": recipeTagID,
	}).Debug("RecipeTagExists called")

	return c.querier.RecipeTagExists(ctx, recipeID, recipeTagID)
}

// GetRecipeTag fetches a recipe tag from the database.
func (c *Client) GetRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) (*models.RecipeTag, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeTag")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeTagIDToSpan(span, recipeTagID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":     recipeID,
		"recipe_tag_id": recipeTagID,
	}).Debug("GetRecipeTag called")

	return c.querier.GetRecipeTag(ctx, recipeID, recipeTagID)
}

// GetAllRecipeTagsCount fetches the count of recipe tags from the database that meet a particular filter.
func (c *Client) GetAllRecipeTagsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeTagsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeTagsCount called")

	return c.querier.GetAllRecipeTagsCount(ctx)
}

// GetRecipeTags fetches a list of recipe tags from the database that meet a particular filter.
func (c *Client) GetRecipeTags(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeTagList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeTags")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
	}).Debug("GetRecipeTags called")

	recipeTagList, err := c.querier.GetRecipeTags(ctx, recipeID, filter)

	return recipeTagList, err
}

// CreateRecipeTag creates a recipe tag in the database.
func (c *Client) CreateRecipeTag(ctx context.Context, input *models.RecipeTagCreationInput) (*models.RecipeTag, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeTag")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeTag called")

	return c.querier.CreateRecipeTag(ctx, input)
}

// UpdateRecipeTag updates a particular recipe tag. Note that UpdateRecipeTag expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeTag(ctx context.Context, updated *models.RecipeTag) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeTag")
	defer span.End()

	tracing.AttachRecipeTagIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_tag_id", updated.ID).Debug("UpdateRecipeTag called")

	return c.querier.UpdateRecipeTag(ctx, updated)
}

// ArchiveRecipeTag archives a recipe tag from the database by its ID.
func (c *Client) ArchiveRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeTag")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeTagIDToSpan(span, recipeTagID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_tag_id": recipeTagID,
		"recipe_id":     recipeID,
	}).Debug("ArchiveRecipeTag called")

	return c.querier.ArchiveRecipeTag(ctx, recipeID, recipeTagID)
}

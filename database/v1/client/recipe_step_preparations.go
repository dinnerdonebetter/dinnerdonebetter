package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeStepPreparationDataManager = (*Client)(nil)

// RecipeStepPreparationExists fetches whether or not a recipe step preparation exists from the database.
func (c *Client) RecipeStepPreparationExists(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepPreparationExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepPreparationIDToSpan(span, recipeStepPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":                  recipeID,
		"recipe_step_id":             recipeStepID,
		"recipe_step_preparation_id": recipeStepPreparationID,
	}).Debug("RecipeStepPreparationExists called")

	return c.querier.RecipeStepPreparationExists(ctx, recipeID, recipeStepID, recipeStepPreparationID)
}

// GetRecipeStepPreparation fetches a recipe step preparation from the database.
func (c *Client) GetRecipeStepPreparation(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (*models.RecipeStepPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepPreparation")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepPreparationIDToSpan(span, recipeStepPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":                  recipeID,
		"recipe_step_id":             recipeStepID,
		"recipe_step_preparation_id": recipeStepPreparationID,
	}).Debug("GetRecipeStepPreparation called")

	return c.querier.GetRecipeStepPreparation(ctx, recipeID, recipeStepID, recipeStepPreparationID)
}

// GetAllRecipeStepPreparationsCount fetches the count of recipe step preparations from the database that meet a particular filter.
func (c *Client) GetAllRecipeStepPreparationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepPreparationsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepPreparationsCount called")

	return c.querier.GetAllRecipeStepPreparationsCount(ctx)
}

// GetRecipeStepPreparations fetches a list of recipe step preparations from the database that meet a particular filter.
func (c *Client) GetRecipeStepPreparations(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepPreparationList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepPreparations")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":      recipeID,
		"recipe_step_id": recipeStepID,
	}).Debug("GetRecipeStepPreparations called")

	recipeStepPreparationList, err := c.querier.GetRecipeStepPreparations(ctx, recipeID, recipeStepID, filter)

	return recipeStepPreparationList, err
}

// CreateRecipeStepPreparation creates a recipe step preparation in the database.
func (c *Client) CreateRecipeStepPreparation(ctx context.Context, input *models.RecipeStepPreparationCreationInput) (*models.RecipeStepPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepPreparation")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepPreparation called")

	return c.querier.CreateRecipeStepPreparation(ctx, input)
}

// UpdateRecipeStepPreparation updates a particular recipe step preparation. Note that UpdateRecipeStepPreparation expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepPreparation(ctx context.Context, updated *models.RecipeStepPreparation) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepPreparation")
	defer span.End()

	tracing.AttachRecipeStepPreparationIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_step_preparation_id", updated.ID).Debug("UpdateRecipeStepPreparation called")

	return c.querier.UpdateRecipeStepPreparation(ctx, updated)
}

// ArchiveRecipeStepPreparation archives a recipe step preparation from the database by its ID.
func (c *Client) ArchiveRecipeStepPreparation(ctx context.Context, recipeStepID, recipeStepPreparationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepPreparation")
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepPreparationIDToSpan(span, recipeStepPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_preparation_id": recipeStepPreparationID,
		"recipe_step_id":             recipeStepID,
	}).Debug("ArchiveRecipeStepPreparation called")

	return c.querier.ArchiveRecipeStepPreparation(ctx, recipeStepID, recipeStepPreparationID)
}

package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeIterationDataManager = (*Client)(nil)

// RecipeIterationExists fetches whether or not a recipe iteration exists from the database.
func (c *Client) RecipeIterationExists(ctx context.Context, recipeID, recipeIterationID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeIterationExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":           recipeID,
		"recipe_iteration_id": recipeIterationID,
	}).Debug("RecipeIterationExists called")

	return c.querier.RecipeIterationExists(ctx, recipeID, recipeIterationID)
}

// GetRecipeIteration fetches a recipe iteration from the database.
func (c *Client) GetRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) (*models.RecipeIteration, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeIteration")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":           recipeID,
		"recipe_iteration_id": recipeIterationID,
	}).Debug("GetRecipeIteration called")

	return c.querier.GetRecipeIteration(ctx, recipeID, recipeIterationID)
}

// GetAllRecipeIterationsCount fetches the count of recipe iterations from the database that meet a particular filter.
func (c *Client) GetAllRecipeIterationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeIterationsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeIterationsCount called")

	return c.querier.GetAllRecipeIterationsCount(ctx)
}

// GetRecipeIterations fetches a list of recipe iterations from the database that meet a particular filter.
func (c *Client) GetRecipeIterations(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeIterationList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeIterations")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
	}).Debug("GetRecipeIterations called")

	recipeIterationList, err := c.querier.GetRecipeIterations(ctx, recipeID, filter)

	return recipeIterationList, err
}

// CreateRecipeIteration creates a recipe iteration in the database.
func (c *Client) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (*models.RecipeIteration, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeIteration")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeIteration called")

	return c.querier.CreateRecipeIteration(ctx, input)
}

// UpdateRecipeIteration updates a particular recipe iteration. Note that UpdateRecipeIteration expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeIteration(ctx context.Context, updated *models.RecipeIteration) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeIteration")
	defer span.End()

	tracing.AttachRecipeIterationIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_iteration_id", updated.ID).Debug("UpdateRecipeIteration called")

	return c.querier.UpdateRecipeIteration(ctx, updated)
}

// ArchiveRecipeIteration archives a recipe iteration from the database by its ID.
func (c *Client) ArchiveRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeIteration")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_iteration_id": recipeIterationID,
		"recipe_id":           recipeID,
	}).Debug("ArchiveRecipeIteration called")

	return c.querier.ArchiveRecipeIteration(ctx, recipeID, recipeIterationID)
}

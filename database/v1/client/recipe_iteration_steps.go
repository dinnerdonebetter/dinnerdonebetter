package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeIterationStepDataManager = (*Client)(nil)

// RecipeIterationStepExists fetches whether or not a recipe iteration step exists from the database.
func (c *Client) RecipeIterationStepExists(ctx context.Context, recipeID, recipeIterationStepID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeIterationStepExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationStepIDToSpan(span, recipeIterationStepID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":                recipeID,
		"recipe_iteration_step_id": recipeIterationStepID,
	}).Debug("RecipeIterationStepExists called")

	return c.querier.RecipeIterationStepExists(ctx, recipeID, recipeIterationStepID)
}

// GetRecipeIterationStep fetches a recipe iteration step from the database.
func (c *Client) GetRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) (*models.RecipeIterationStep, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeIterationStep")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationStepIDToSpan(span, recipeIterationStepID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":                recipeID,
		"recipe_iteration_step_id": recipeIterationStepID,
	}).Debug("GetRecipeIterationStep called")

	return c.querier.GetRecipeIterationStep(ctx, recipeID, recipeIterationStepID)
}

// GetAllRecipeIterationStepsCount fetches the count of recipe iteration steps from the database that meet a particular filter.
func (c *Client) GetAllRecipeIterationStepsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeIterationStepsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeIterationStepsCount called")

	return c.querier.GetAllRecipeIterationStepsCount(ctx)
}

// GetRecipeIterationSteps fetches a list of recipe iteration steps from the database that meet a particular filter.
func (c *Client) GetRecipeIterationSteps(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeIterationStepList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeIterationSteps")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
	}).Debug("GetRecipeIterationSteps called")

	recipeIterationStepList, err := c.querier.GetRecipeIterationSteps(ctx, recipeID, filter)

	return recipeIterationStepList, err
}

// CreateRecipeIterationStep creates a recipe iteration step in the database.
func (c *Client) CreateRecipeIterationStep(ctx context.Context, input *models.RecipeIterationStepCreationInput) (*models.RecipeIterationStep, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeIterationStep")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeIterationStep called")

	return c.querier.CreateRecipeIterationStep(ctx, input)
}

// UpdateRecipeIterationStep updates a particular recipe iteration step. Note that UpdateRecipeIterationStep expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeIterationStep(ctx context.Context, updated *models.RecipeIterationStep) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeIterationStep")
	defer span.End()

	tracing.AttachRecipeIterationStepIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_iteration_step_id", updated.ID).Debug("UpdateRecipeIterationStep called")

	return c.querier.UpdateRecipeIterationStep(ctx, updated)
}

// ArchiveRecipeIterationStep archives a recipe iteration step from the database by its ID.
func (c *Client) ArchiveRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeIterationStep")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationStepIDToSpan(span, recipeIterationStepID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_iteration_step_id": recipeIterationStepID,
		"recipe_id":                recipeID,
	}).Debug("ArchiveRecipeIterationStep called")

	return c.querier.ArchiveRecipeIterationStep(ctx, recipeID, recipeIterationStepID)
}

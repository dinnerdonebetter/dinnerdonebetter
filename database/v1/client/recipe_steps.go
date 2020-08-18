package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeStepDataManager = (*Client)(nil)

// RecipeStepExists fetches whether or not a recipe step exists from the database.
func (c *Client) RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":      recipeID,
		"recipe_step_id": recipeStepID,
	}).Debug("RecipeStepExists called")

	return c.querier.RecipeStepExists(ctx, recipeID, recipeStepID)
}

// GetRecipeStep fetches a recipe step from the database.
func (c *Client) GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*models.RecipeStep, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStep")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":      recipeID,
		"recipe_step_id": recipeStepID,
	}).Debug("GetRecipeStep called")

	return c.querier.GetRecipeStep(ctx, recipeID, recipeStepID)
}

// GetAllRecipeStepsCount fetches the count of recipe steps from the database that meet a particular filter.
func (c *Client) GetAllRecipeStepsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepsCount called")

	return c.querier.GetAllRecipeStepsCount(ctx)
}

// GetAllRecipeSteps fetches a list of all recipe steps in the database.
func (c *Client) GetAllRecipeSteps(ctx context.Context, results chan []models.RecipeStep) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeSteps")
	defer span.End()

	c.logger.Debug("GetAllRecipeSteps called")

	return c.querier.GetAllRecipeSteps(ctx, results)
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter.
func (c *Client) GetRecipeSteps(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeStepList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeSteps")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id": recipeID,
	}).Debug("GetRecipeSteps called")

	recipeStepList, err := c.querier.GetRecipeSteps(ctx, recipeID, filter)

	return recipeStepList, err
}

// GetRecipeStepsWithIDs fetches recipe steps from the database within a given set of IDs.
func (c *Client) GetRecipeStepsWithIDs(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) ([]models.RecipeStep, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetRecipeStepsWithIDs called")

	recipeStepList, err := c.querier.GetRecipeStepsWithIDs(ctx, recipeID, limit, ids)

	return recipeStepList, err
}

// CreateRecipeStep creates a recipe step in the database.
func (c *Client) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (*models.RecipeStep, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStep")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStep called")

	return c.querier.CreateRecipeStep(ctx, input)
}

// UpdateRecipeStep updates a particular recipe step. Note that UpdateRecipeStep expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStep(ctx context.Context, updated *models.RecipeStep) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStep")
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_step_id", updated.ID).Debug("UpdateRecipeStep called")

	return c.querier.UpdateRecipeStep(ctx, updated)
}

// ArchiveRecipeStep archives a recipe step from the database by its ID.
func (c *Client) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStep")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_id": recipeStepID,
		"recipe_id":      recipeID,
	}).Debug("ArchiveRecipeStep called")

	return c.querier.ArchiveRecipeStep(ctx, recipeID, recipeStepID)
}

package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.RecipeStepDataManager = (*Client)(nil)

// attachRecipeStepIDToSpan provides a consistent way to attach a recipe step's ID to a span
func attachRecipeStepIDToSpan(span *trace.Span, recipeStepID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_step_id", strconv.FormatUint(recipeStepID, 10)))
	}
}

// GetRecipeStep fetches a recipe step from the database
func (c *Client) GetRecipeStep(ctx context.Context, recipeStepID, userID uint64) (*models.RecipeStep, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStep")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepIDToSpan(span, recipeStepID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_id": recipeStepID,
		"user_id":        userID,
	}).Debug("GetRecipeStep called")

	return c.querier.GetRecipeStep(ctx, recipeStepID, userID)
}

// GetRecipeStepCount fetches the count of recipe steps from the database that meet a particular filter
func (c *Client) GetRecipeStepCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepCount called")

	return c.querier.GetRecipeStepCount(ctx, filter, userID)
}

// GetAllRecipeStepsCount fetches the count of recipe steps from the database that meet a particular filter
func (c *Client) GetAllRecipeStepsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepsCount called")

	return c.querier.GetAllRecipeStepsCount(ctx)
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter
func (c *Client) GetRecipeSteps(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepList, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeSteps")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeSteps called")

	recipeStepList, err := c.querier.GetRecipeSteps(ctx, filter, userID)

	return recipeStepList, err
}

// GetAllRecipeStepsForUser fetches a list of recipe steps from the database that meet a particular filter
func (c *Client) GetAllRecipeStepsForUser(ctx context.Context, userID uint64) ([]models.RecipeStep, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllRecipeStepsForUser called")

	recipeStepList, err := c.querier.GetAllRecipeStepsForUser(ctx, userID)

	return recipeStepList, err
}

// CreateRecipeStep creates a recipe step in the database
func (c *Client) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (*models.RecipeStep, error) {
	ctx, span := trace.StartSpan(ctx, "CreateRecipeStep")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStep called")

	return c.querier.CreateRecipeStep(ctx, input)
}

// UpdateRecipeStep updates a particular recipe step. Note that UpdateRecipeStep expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStep(ctx context.Context, input *models.RecipeStep) error {
	ctx, span := trace.StartSpan(ctx, "UpdateRecipeStep")
	defer span.End()

	attachRecipeStepIDToSpan(span, input.ID)
	c.logger.WithValue("recipe_step_id", input.ID).Debug("UpdateRecipeStep called")

	return c.querier.UpdateRecipeStep(ctx, input)
}

// ArchiveRecipeStep archives a recipe step from the database by its ID
func (c *Client) ArchiveRecipeStep(ctx context.Context, recipeStepID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveRecipeStep")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepIDToSpan(span, recipeStepID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_id": recipeStepID,
		"user_id":        userID,
	}).Debug("ArchiveRecipeStep called")

	return c.querier.ArchiveRecipeStep(ctx, recipeStepID, userID)
}

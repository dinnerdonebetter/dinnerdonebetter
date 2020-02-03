package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.RecipeIterationDataManager = (*Client)(nil)

// attachRecipeIterationIDToSpan provides a consistent way to attach a recipe iteration's ID to a span
func attachRecipeIterationIDToSpan(span *trace.Span, recipeIterationID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_iteration_id", strconv.FormatUint(recipeIterationID, 10)))
	}
}

// GetRecipeIteration fetches a recipe iteration from the database
func (c *Client) GetRecipeIteration(ctx context.Context, recipeIterationID, userID uint64) (*models.RecipeIteration, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeIteration")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeIterationIDToSpan(span, recipeIterationID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_iteration_id": recipeIterationID,
		"user_id":             userID,
	}).Debug("GetRecipeIteration called")

	return c.querier.GetRecipeIteration(ctx, recipeIterationID, userID)
}

// GetRecipeIterationCount fetches the count of recipe iterations from the database that meet a particular filter
func (c *Client) GetRecipeIterationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeIterationCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeIterationCount called")

	return c.querier.GetRecipeIterationCount(ctx, filter, userID)
}

// GetAllRecipeIterationsCount fetches the count of recipe iterations from the database that meet a particular filter
func (c *Client) GetAllRecipeIterationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeIterationsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeIterationsCount called")

	return c.querier.GetAllRecipeIterationsCount(ctx)
}

// GetRecipeIterations fetches a list of recipe iterations from the database that meet a particular filter
func (c *Client) GetRecipeIterations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeIterationList, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeIterations")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeIterations called")

	recipeIterationList, err := c.querier.GetRecipeIterations(ctx, filter, userID)

	return recipeIterationList, err
}

// GetAllRecipeIterationsForUser fetches a list of recipe iterations from the database that meet a particular filter
func (c *Client) GetAllRecipeIterationsForUser(ctx context.Context, userID uint64) ([]models.RecipeIteration, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeIterationsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllRecipeIterationsForUser called")

	recipeIterationList, err := c.querier.GetAllRecipeIterationsForUser(ctx, userID)

	return recipeIterationList, err
}

// CreateRecipeIteration creates a recipe iteration in the database
func (c *Client) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (*models.RecipeIteration, error) {
	ctx, span := trace.StartSpan(ctx, "CreateRecipeIteration")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeIteration called")

	return c.querier.CreateRecipeIteration(ctx, input)
}

// UpdateRecipeIteration updates a particular recipe iteration. Note that UpdateRecipeIteration expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeIteration(ctx context.Context, input *models.RecipeIteration) error {
	ctx, span := trace.StartSpan(ctx, "UpdateRecipeIteration")
	defer span.End()

	attachRecipeIterationIDToSpan(span, input.ID)
	c.logger.WithValue("recipe_iteration_id", input.ID).Debug("UpdateRecipeIteration called")

	return c.querier.UpdateRecipeIteration(ctx, input)
}

// ArchiveRecipeIteration archives a recipe iteration from the database by its ID
func (c *Client) ArchiveRecipeIteration(ctx context.Context, recipeIterationID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveRecipeIteration")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeIterationIDToSpan(span, recipeIterationID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_iteration_id": recipeIterationID,
		"user_id":             userID,
	}).Debug("ArchiveRecipeIteration called")

	return c.querier.ArchiveRecipeIteration(ctx, recipeIterationID, userID)
}

package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeStepEventDataManager = (*Client)(nil)

// RecipeStepEventExists fetches whether or not a recipe step event exists from the database.
func (c *Client) RecipeStepEventExists(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepEventExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepEventIDToSpan(span, recipeStepEventID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":            recipeID,
		"recipe_step_id":       recipeStepID,
		"recipe_step_event_id": recipeStepEventID,
	}).Debug("RecipeStepEventExists called")

	return c.querier.RecipeStepEventExists(ctx, recipeID, recipeStepID, recipeStepEventID)
}

// GetRecipeStepEvent fetches a recipe step event from the database.
func (c *Client) GetRecipeStepEvent(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (*models.RecipeStepEvent, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepEvent")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepEventIDToSpan(span, recipeStepEventID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":            recipeID,
		"recipe_step_id":       recipeStepID,
		"recipe_step_event_id": recipeStepEventID,
	}).Debug("GetRecipeStepEvent called")

	return c.querier.GetRecipeStepEvent(ctx, recipeID, recipeStepID, recipeStepEventID)
}

// GetAllRecipeStepEventsCount fetches the count of recipe step events from the database that meet a particular filter.
func (c *Client) GetAllRecipeStepEventsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepEventsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepEventsCount called")

	return c.querier.GetAllRecipeStepEventsCount(ctx)
}

// GetAllRecipeStepEvents fetches a list of all recipe step events in the database.
func (c *Client) GetAllRecipeStepEvents(ctx context.Context, results chan []models.RecipeStepEvent) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepEvents")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepEvents called")

	return c.querier.GetAllRecipeStepEvents(ctx, results)
}

// GetRecipeStepEvents fetches a list of recipe step events from the database that meet a particular filter.
func (c *Client) GetRecipeStepEvents(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepEventList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepEvents")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":      recipeID,
		"recipe_step_id": recipeStepID,
	}).Debug("GetRecipeStepEvents called")

	recipeStepEventList, err := c.querier.GetRecipeStepEvents(ctx, recipeID, recipeStepID, filter)

	return recipeStepEventList, err
}

// GetRecipeStepEventsWithIDs fetches recipe step events from the database within a given set of IDs.
func (c *Client) GetRecipeStepEventsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepEvent, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepEventsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetRecipeStepEventsWithIDs called")

	recipeStepEventList, err := c.querier.GetRecipeStepEventsWithIDs(ctx, recipeID, recipeStepID, limit, ids)

	return recipeStepEventList, err
}

// CreateRecipeStepEvent creates a recipe step event in the database.
func (c *Client) CreateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEventCreationInput) (*models.RecipeStepEvent, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepEvent")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepEvent called")

	return c.querier.CreateRecipeStepEvent(ctx, input)
}

// UpdateRecipeStepEvent updates a particular recipe step event. Note that UpdateRecipeStepEvent expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepEvent(ctx context.Context, updated *models.RecipeStepEvent) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepEvent")
	defer span.End()

	tracing.AttachRecipeStepEventIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_step_event_id", updated.ID).Debug("UpdateRecipeStepEvent called")

	return c.querier.UpdateRecipeStepEvent(ctx, updated)
}

// ArchiveRecipeStepEvent archives a recipe step event from the database by its ID.
func (c *Client) ArchiveRecipeStepEvent(ctx context.Context, recipeStepID, recipeStepEventID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepEvent")
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepEventIDToSpan(span, recipeStepEventID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_event_id": recipeStepEventID,
		"recipe_step_id":       recipeStepID,
	}).Debug("ArchiveRecipeStepEvent called")

	return c.querier.ArchiveRecipeStepEvent(ctx, recipeStepID, recipeStepEventID)
}

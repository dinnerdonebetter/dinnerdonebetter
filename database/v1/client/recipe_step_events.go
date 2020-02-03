package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.RecipeStepEventDataManager = (*Client)(nil)

// attachRecipeStepEventIDToSpan provides a consistent way to attach a recipe step event's ID to a span
func attachRecipeStepEventIDToSpan(span *trace.Span, recipeStepEventID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_step_event_id", strconv.FormatUint(recipeStepEventID, 10)))
	}
}

// GetRecipeStepEvent fetches a recipe step event from the database
func (c *Client) GetRecipeStepEvent(ctx context.Context, recipeStepEventID, userID uint64) (*models.RecipeStepEvent, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepEvent")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepEventIDToSpan(span, recipeStepEventID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_event_id": recipeStepEventID,
		"user_id":              userID,
	}).Debug("GetRecipeStepEvent called")

	return c.querier.GetRecipeStepEvent(ctx, recipeStepEventID, userID)
}

// GetRecipeStepEventCount fetches the count of recipe step events from the database that meet a particular filter
func (c *Client) GetRecipeStepEventCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepEventCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepEventCount called")

	return c.querier.GetRecipeStepEventCount(ctx, filter, userID)
}

// GetAllRecipeStepEventsCount fetches the count of recipe step events from the database that meet a particular filter
func (c *Client) GetAllRecipeStepEventsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepEventsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepEventsCount called")

	return c.querier.GetAllRecipeStepEventsCount(ctx)
}

// GetRecipeStepEvents fetches a list of recipe step events from the database that meet a particular filter
func (c *Client) GetRecipeStepEvents(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepEventList, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepEvents")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepEvents called")

	recipeStepEventList, err := c.querier.GetRecipeStepEvents(ctx, filter, userID)

	return recipeStepEventList, err
}

// GetAllRecipeStepEventsForUser fetches a list of recipe step events from the database that meet a particular filter
func (c *Client) GetAllRecipeStepEventsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepEvent, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepEventsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllRecipeStepEventsForUser called")

	recipeStepEventList, err := c.querier.GetAllRecipeStepEventsForUser(ctx, userID)

	return recipeStepEventList, err
}

// CreateRecipeStepEvent creates a recipe step event in the database
func (c *Client) CreateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEventCreationInput) (*models.RecipeStepEvent, error) {
	ctx, span := trace.StartSpan(ctx, "CreateRecipeStepEvent")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepEvent called")

	return c.querier.CreateRecipeStepEvent(ctx, input)
}

// UpdateRecipeStepEvent updates a particular recipe step event. Note that UpdateRecipeStepEvent expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEvent) error {
	ctx, span := trace.StartSpan(ctx, "UpdateRecipeStepEvent")
	defer span.End()

	attachRecipeStepEventIDToSpan(span, input.ID)
	c.logger.WithValue("recipe_step_event_id", input.ID).Debug("UpdateRecipeStepEvent called")

	return c.querier.UpdateRecipeStepEvent(ctx, input)
}

// ArchiveRecipeStepEvent archives a recipe step event from the database by its ID
func (c *Client) ArchiveRecipeStepEvent(ctx context.Context, recipeStepEventID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveRecipeStepEvent")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepEventIDToSpan(span, recipeStepEventID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_event_id": recipeStepEventID,
		"user_id":              userID,
	}).Debug("ArchiveRecipeStepEvent called")

	return c.querier.ArchiveRecipeStepEvent(ctx, recipeStepEventID, userID)
}

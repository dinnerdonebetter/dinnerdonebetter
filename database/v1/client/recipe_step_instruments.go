package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.RecipeStepInstrumentDataManager = (*Client)(nil)

// attachRecipeStepInstrumentIDToSpan provides a consistent way to attach a recipe step instrument's ID to a span
func attachRecipeStepInstrumentIDToSpan(span *trace.Span, recipeStepInstrumentID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_step_instrument_id", strconv.FormatUint(recipeStepInstrumentID, 10)))
	}
}

// GetRecipeStepInstrument fetches a recipe step instrument from the database
func (c *Client) GetRecipeStepInstrument(ctx context.Context, recipeStepInstrumentID, userID uint64) (*models.RecipeStepInstrument, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepInstrument")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_instrument_id": recipeStepInstrumentID,
		"user_id":                   userID,
	}).Debug("GetRecipeStepInstrument called")

	return c.querier.GetRecipeStepInstrument(ctx, recipeStepInstrumentID, userID)
}

// GetRecipeStepInstrumentCount fetches the count of recipe step instruments from the database that meet a particular filter
func (c *Client) GetRecipeStepInstrumentCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepInstrumentCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepInstrumentCount called")

	return c.querier.GetRecipeStepInstrumentCount(ctx, filter, userID)
}

// GetAllRecipeStepInstrumentsCount fetches the count of recipe step instruments from the database that meet a particular filter
func (c *Client) GetAllRecipeStepInstrumentsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepInstrumentsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepInstrumentsCount called")

	return c.querier.GetAllRecipeStepInstrumentsCount(ctx)
}

// GetRecipeStepInstruments fetches a list of recipe step instruments from the database that meet a particular filter
func (c *Client) GetRecipeStepInstruments(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepInstrumentList, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepInstruments")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepInstruments called")

	recipeStepInstrumentList, err := c.querier.GetRecipeStepInstruments(ctx, filter, userID)

	return recipeStepInstrumentList, err
}

// GetAllRecipeStepInstrumentsForUser fetches a list of recipe step instruments from the database that meet a particular filter
func (c *Client) GetAllRecipeStepInstrumentsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepInstrument, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepInstrumentsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllRecipeStepInstrumentsForUser called")

	recipeStepInstrumentList, err := c.querier.GetAllRecipeStepInstrumentsForUser(ctx, userID)

	return recipeStepInstrumentList, err
}

// CreateRecipeStepInstrument creates a recipe step instrument in the database
func (c *Client) CreateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrumentCreationInput) (*models.RecipeStepInstrument, error) {
	ctx, span := trace.StartSpan(ctx, "CreateRecipeStepInstrument")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepInstrument called")

	return c.querier.CreateRecipeStepInstrument(ctx, input)
}

// UpdateRecipeStepInstrument updates a particular recipe step instrument. Note that UpdateRecipeStepInstrument expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrument) error {
	ctx, span := trace.StartSpan(ctx, "UpdateRecipeStepInstrument")
	defer span.End()

	attachRecipeStepInstrumentIDToSpan(span, input.ID)
	c.logger.WithValue("recipe_step_instrument_id", input.ID).Debug("UpdateRecipeStepInstrument called")

	return c.querier.UpdateRecipeStepInstrument(ctx, input)
}

// ArchiveRecipeStepInstrument archives a recipe step instrument from the database by its ID
func (c *Client) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepInstrumentID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveRecipeStepInstrument")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_instrument_id": recipeStepInstrumentID,
		"user_id":                   userID,
	}).Debug("ArchiveRecipeStepInstrument called")

	return c.querier.ArchiveRecipeStepInstrument(ctx, recipeStepInstrumentID, userID)
}

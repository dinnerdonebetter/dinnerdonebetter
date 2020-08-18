package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeStepInstrumentDataManager = (*Client)(nil)

// RecipeStepInstrumentExists fetches whether or not a recipe step instrument exists from the database.
func (c *Client) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepInstrumentExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":                 recipeID,
		"recipe_step_id":            recipeStepID,
		"recipe_step_instrument_id": recipeStepInstrumentID,
	}).Debug("RecipeStepInstrumentExists called")

	return c.querier.RecipeStepInstrumentExists(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
}

// GetRecipeStepInstrument fetches a recipe step instrument from the database.
func (c *Client) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (*models.RecipeStepInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepInstrument")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":                 recipeID,
		"recipe_step_id":            recipeStepID,
		"recipe_step_instrument_id": recipeStepInstrumentID,
	}).Debug("GetRecipeStepInstrument called")

	return c.querier.GetRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
}

// GetAllRecipeStepInstrumentsCount fetches the count of recipe step instruments from the database that meet a particular filter.
func (c *Client) GetAllRecipeStepInstrumentsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepInstrumentsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepInstrumentsCount called")

	return c.querier.GetAllRecipeStepInstrumentsCount(ctx)
}

// GetAllRecipeStepInstruments fetches a list of all recipe step instruments in the database.
func (c *Client) GetAllRecipeStepInstruments(ctx context.Context, results chan []models.RecipeStepInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepInstruments")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepInstruments called")

	return c.querier.GetAllRecipeStepInstruments(ctx, results)
}

// GetRecipeStepInstruments fetches a list of recipe step instruments from the database that meet a particular filter.
func (c *Client) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepInstrumentList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepInstruments")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":      recipeID,
		"recipe_step_id": recipeStepID,
	}).Debug("GetRecipeStepInstruments called")

	recipeStepInstrumentList, err := c.querier.GetRecipeStepInstruments(ctx, recipeID, recipeStepID, filter)

	return recipeStepInstrumentList, err
}

// GetRecipeStepInstrumentsWithIDs fetches recipe step instruments from the database within a given set of IDs.
func (c *Client) GetRecipeStepInstrumentsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepInstrumentsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetRecipeStepInstrumentsWithIDs called")

	recipeStepInstrumentList, err := c.querier.GetRecipeStepInstrumentsWithIDs(ctx, recipeID, recipeStepID, limit, ids)

	return recipeStepInstrumentList, err
}

// CreateRecipeStepInstrument creates a recipe step instrument in the database.
func (c *Client) CreateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrumentCreationInput) (*models.RecipeStepInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepInstrument")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepInstrument called")

	return c.querier.CreateRecipeStepInstrument(ctx, input)
}

// UpdateRecipeStepInstrument updates a particular recipe step instrument. Note that UpdateRecipeStepInstrument expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepInstrument(ctx context.Context, updated *models.RecipeStepInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepInstrument")
	defer span.End()

	tracing.AttachRecipeStepInstrumentIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_step_instrument_id", updated.ID).Debug("UpdateRecipeStepInstrument called")

	return c.querier.UpdateRecipeStepInstrument(ctx, updated)
}

// ArchiveRecipeStepInstrument archives a recipe step instrument from the database by its ID.
func (c *Client) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepInstrument")
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_instrument_id": recipeStepInstrumentID,
		"recipe_step_id":            recipeStepID,
	}).Debug("ArchiveRecipeStepInstrument called")

	return c.querier.ArchiveRecipeStepInstrument(ctx, recipeStepID, recipeStepInstrumentID)
}

package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.ValidIngredientDataManager = (*Client)(nil)

// ValidIngredientExists fetches whether or not a valid ingredient exists from the database.
func (c *Client) ValidIngredientExists(ctx context.Context, validIngredientID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "ValidIngredientExists")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id": validIngredientID,
	}).Debug("ValidIngredientExists called")

	return c.querier.ValidIngredientExists(ctx, validIngredientID)
}

// GetValidIngredient fetches a valid ingredient from the database.
func (c *Client) GetValidIngredient(ctx context.Context, validIngredientID uint64) (*models.ValidIngredient, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredient")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id": validIngredientID,
	}).Debug("GetValidIngredient called")

	return c.querier.GetValidIngredient(ctx, validIngredientID)
}

// GetAllValidIngredientsCount fetches the count of valid ingredients from the database that meet a particular filter.
func (c *Client) GetAllValidIngredientsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidIngredientsCount")
	defer span.End()

	c.logger.Debug("GetAllValidIngredientsCount called")

	return c.querier.GetAllValidIngredientsCount(ctx)
}

// GetAllValidIngredients fetches a list of all valid ingredients in the database.
func (c *Client) GetAllValidIngredients(ctx context.Context, results chan []models.ValidIngredient) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidIngredients")
	defer span.End()

	c.logger.Debug("GetAllValidIngredients called")

	return c.querier.GetAllValidIngredients(ctx, results)
}

// GetValidIngredients fetches a list of valid ingredients from the database that meet a particular filter.
func (c *Client) GetValidIngredients(ctx context.Context, filter *models.QueryFilter) (*models.ValidIngredientList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredients")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetValidIngredients called")

	validIngredientList, err := c.querier.GetValidIngredients(ctx, filter)

	return validIngredientList, err
}

// GetValidIngredientsWithIDs fetches valid ingredients from the database within a given set of IDs.
func (c *Client) GetValidIngredientsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidIngredient, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetValidIngredientsWithIDs called")

	validIngredientList, err := c.querier.GetValidIngredientsWithIDs(ctx, limit, ids)

	return validIngredientList, err
}

// CreateValidIngredient creates a valid ingredient in the database.
func (c *Client) CreateValidIngredient(ctx context.Context, input *models.ValidIngredientCreationInput) (*models.ValidIngredient, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidIngredient")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateValidIngredient called")

	return c.querier.CreateValidIngredient(ctx, input)
}

// UpdateValidIngredient updates a particular valid ingredient. Note that UpdateValidIngredient expects the
// provided input to have a valid ID.
func (c *Client) UpdateValidIngredient(ctx context.Context, updated *models.ValidIngredient) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidIngredient")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, updated.ID)
	c.logger.WithValue("valid_ingredient_id", updated.ID).Debug("UpdateValidIngredient called")

	return c.querier.UpdateValidIngredient(ctx, updated)
}

// ArchiveValidIngredient archives a valid ingredient from the database by its ID.
func (c *Client) ArchiveValidIngredient(ctx context.Context, validIngredientID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidIngredient")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id": validIngredientID,
	}).Debug("ArchiveValidIngredient called")

	return c.querier.ArchiveValidIngredient(ctx, validIngredientID)
}

package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.ValidIngredientPreparationDataManager = (*Client)(nil)

// ValidIngredientPreparationExists fetches whether or not a valid ingredient preparation exists from the database.
func (c *Client) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "ValidIngredientPreparationExists")
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_preparation_id": validIngredientPreparationID,
	}).Debug("ValidIngredientPreparationExists called")

	return c.querier.ValidIngredientPreparationExists(ctx, validIngredientPreparationID)
}

// GetValidIngredientPreparation fetches a valid ingredient preparation from the database.
func (c *Client) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) (*models.ValidIngredientPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientPreparation")
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_preparation_id": validIngredientPreparationID,
	}).Debug("GetValidIngredientPreparation called")

	return c.querier.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
}

// GetAllValidIngredientPreparationsCount fetches the count of valid ingredient preparations from the database that meet a particular filter.
func (c *Client) GetAllValidIngredientPreparationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidIngredientPreparationsCount")
	defer span.End()

	c.logger.Debug("GetAllValidIngredientPreparationsCount called")

	return c.querier.GetAllValidIngredientPreparationsCount(ctx)
}

// GetAllValidIngredientPreparations fetches a list of all valid ingredient preparations in the database.
func (c *Client) GetAllValidIngredientPreparations(ctx context.Context, results chan []models.ValidIngredientPreparation) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidIngredientPreparations")
	defer span.End()

	c.logger.Debug("GetAllValidIngredientPreparations called")

	return c.querier.GetAllValidIngredientPreparations(ctx, results)
}

// GetValidIngredientPreparations fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (c *Client) GetValidIngredientPreparations(ctx context.Context, filter *models.QueryFilter) (*models.ValidIngredientPreparationList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientPreparations")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetValidIngredientPreparations called")

	validIngredientPreparationList, err := c.querier.GetValidIngredientPreparations(ctx, filter)

	return validIngredientPreparationList, err
}

// GetValidIngredientPreparationsWithIDs fetches valid ingredient preparations from the database within a given set of IDs.
func (c *Client) GetValidIngredientPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidIngredientPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientPreparationsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetValidIngredientPreparationsWithIDs called")

	validIngredientPreparationList, err := c.querier.GetValidIngredientPreparationsWithIDs(ctx, limit, ids)

	return validIngredientPreparationList, err
}

// CreateValidIngredientPreparation creates a valid ingredient preparation in the database.
func (c *Client) CreateValidIngredientPreparation(ctx context.Context, input *models.ValidIngredientPreparationCreationInput) (*models.ValidIngredientPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidIngredientPreparation")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateValidIngredientPreparation called")

	return c.querier.CreateValidIngredientPreparation(ctx, input)
}

// UpdateValidIngredientPreparation updates a particular valid ingredient preparation. Note that UpdateValidIngredientPreparation expects the
// provided input to have a valid ID.
func (c *Client) UpdateValidIngredientPreparation(ctx context.Context, updated *models.ValidIngredientPreparation) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidIngredientPreparation")
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, updated.ID)
	c.logger.WithValue("valid_ingredient_preparation_id", updated.ID).Debug("UpdateValidIngredientPreparation called")

	return c.querier.UpdateValidIngredientPreparation(ctx, updated)
}

// ArchiveValidIngredientPreparation archives a valid ingredient preparation from the database by its ID.
func (c *Client) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidIngredientPreparation")
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_preparation_id": validIngredientPreparationID,
	}).Debug("ArchiveValidIngredientPreparation called")

	return c.querier.ArchiveValidIngredientPreparation(ctx, validIngredientPreparationID)
}

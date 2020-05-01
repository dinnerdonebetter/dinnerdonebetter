package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.ValidIngredientPreparationDataManager = (*Client)(nil)

// ValidIngredientPreparationExists fetches whether or not a valid ingredient preparation exists from the database.
func (c *Client) ValidIngredientPreparationExists(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "ValidIngredientPreparationExists")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id":             validIngredientID,
		"valid_ingredient_preparation_id": validIngredientPreparationID,
	}).Debug("ValidIngredientPreparationExists called")

	return c.querier.ValidIngredientPreparationExists(ctx, validIngredientID, validIngredientPreparationID)
}

// GetValidIngredientPreparation fetches a valid ingredient preparation from the database.
func (c *Client) GetValidIngredientPreparation(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) (*models.ValidIngredientPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientPreparation")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id":             validIngredientID,
		"valid_ingredient_preparation_id": validIngredientPreparationID,
	}).Debug("GetValidIngredientPreparation called")

	return c.querier.GetValidIngredientPreparation(ctx, validIngredientID, validIngredientPreparationID)
}

// GetAllValidIngredientPreparationsCount fetches the count of valid ingredient preparations from the database that meet a particular filter.
func (c *Client) GetAllValidIngredientPreparationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidIngredientPreparationsCount")
	defer span.End()

	c.logger.Debug("GetAllValidIngredientPreparationsCount called")

	return c.querier.GetAllValidIngredientPreparationsCount(ctx)
}

// GetValidIngredientPreparations fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (c *Client) GetValidIngredientPreparations(ctx context.Context, validIngredientID uint64, filter *models.QueryFilter) (*models.ValidIngredientPreparationList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientPreparations")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id": validIngredientID,
	}).Debug("GetValidIngredientPreparations called")

	validIngredientPreparationList, err := c.querier.GetValidIngredientPreparations(ctx, validIngredientID, filter)

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
func (c *Client) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidIngredientPreparation")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_preparation_id": validIngredientPreparationID,
		"valid_ingredient_id":             validIngredientID,
	}).Debug("ArchiveValidIngredientPreparation called")

	return c.querier.ArchiveValidIngredientPreparation(ctx, validIngredientID, validIngredientPreparationID)
}

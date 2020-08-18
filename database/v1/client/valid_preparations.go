package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.ValidPreparationDataManager = (*Client)(nil)

// ValidPreparationExists fetches whether or not a valid preparation exists from the database.
func (c *Client) ValidPreparationExists(ctx context.Context, validPreparationID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "ValidPreparationExists")
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_preparation_id": validPreparationID,
	}).Debug("ValidPreparationExists called")

	return c.querier.ValidPreparationExists(ctx, validPreparationID)
}

// GetValidPreparation fetches a valid preparation from the database.
func (c *Client) GetValidPreparation(ctx context.Context, validPreparationID uint64) (*models.ValidPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidPreparation")
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_preparation_id": validPreparationID,
	}).Debug("GetValidPreparation called")

	return c.querier.GetValidPreparation(ctx, validPreparationID)
}

// GetAllValidPreparationsCount fetches the count of valid preparations from the database that meet a particular filter.
func (c *Client) GetAllValidPreparationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidPreparationsCount")
	defer span.End()

	c.logger.Debug("GetAllValidPreparationsCount called")

	return c.querier.GetAllValidPreparationsCount(ctx)
}

// GetAllValidPreparations fetches a list of all valid preparations in the database.
func (c *Client) GetAllValidPreparations(ctx context.Context, results chan []models.ValidPreparation) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidPreparations")
	defer span.End()

	c.logger.Debug("GetAllValidPreparations called")

	return c.querier.GetAllValidPreparations(ctx, results)
}

// GetValidPreparations fetches a list of valid preparations from the database that meet a particular filter.
func (c *Client) GetValidPreparations(ctx context.Context, filter *models.QueryFilter) (*models.ValidPreparationList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidPreparations")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetValidPreparations called")

	validPreparationList, err := c.querier.GetValidPreparations(ctx, filter)

	return validPreparationList, err
}

// GetValidPreparationsWithIDs fetches valid preparations from the database within a given set of IDs.
func (c *Client) GetValidPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidPreparationsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetValidPreparationsWithIDs called")

	validPreparationList, err := c.querier.GetValidPreparationsWithIDs(ctx, limit, ids)

	return validPreparationList, err
}

// CreateValidPreparation creates a valid preparation in the database.
func (c *Client) CreateValidPreparation(ctx context.Context, input *models.ValidPreparationCreationInput) (*models.ValidPreparation, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidPreparation")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateValidPreparation called")

	return c.querier.CreateValidPreparation(ctx, input)
}

// UpdateValidPreparation updates a particular valid preparation. Note that UpdateValidPreparation expects the
// provided input to have a valid ID.
func (c *Client) UpdateValidPreparation(ctx context.Context, updated *models.ValidPreparation) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidPreparation")
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, updated.ID)
	c.logger.WithValue("valid_preparation_id", updated.ID).Debug("UpdateValidPreparation called")

	return c.querier.UpdateValidPreparation(ctx, updated)
}

// ArchiveValidPreparation archives a valid preparation from the database by its ID.
func (c *Client) ArchiveValidPreparation(ctx context.Context, validPreparationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidPreparation")
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	c.logger.WithValues(map[string]interface{}{
		"valid_preparation_id": validPreparationID,
	}).Debug("ArchiveValidPreparation called")

	return c.querier.ArchiveValidPreparation(ctx, validPreparationID)
}

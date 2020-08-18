package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RequiredPreparationInstrumentDataManager = (*Client)(nil)

// RequiredPreparationInstrumentExists fetches whether or not a required preparation instrument exists from the database.
func (c *Client) RequiredPreparationInstrumentExists(ctx context.Context, requiredPreparationInstrumentID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RequiredPreparationInstrumentExists")
	defer span.End()

	tracing.AttachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"required_preparation_instrument_id": requiredPreparationInstrumentID,
	}).Debug("RequiredPreparationInstrumentExists called")

	return c.querier.RequiredPreparationInstrumentExists(ctx, requiredPreparationInstrumentID)
}

// GetRequiredPreparationInstrument fetches a required preparation instrument from the database.
func (c *Client) GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) (*models.RequiredPreparationInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRequiredPreparationInstrument")
	defer span.End()

	tracing.AttachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"required_preparation_instrument_id": requiredPreparationInstrumentID,
	}).Debug("GetRequiredPreparationInstrument called")

	return c.querier.GetRequiredPreparationInstrument(ctx, requiredPreparationInstrumentID)
}

// GetAllRequiredPreparationInstrumentsCount fetches the count of required preparation instruments from the database that meet a particular filter.
func (c *Client) GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRequiredPreparationInstrumentsCount")
	defer span.End()

	c.logger.Debug("GetAllRequiredPreparationInstrumentsCount called")

	return c.querier.GetAllRequiredPreparationInstrumentsCount(ctx)
}

// GetAllRequiredPreparationInstruments fetches a list of all required preparation instruments in the database.
func (c *Client) GetAllRequiredPreparationInstruments(ctx context.Context, results chan []models.RequiredPreparationInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllRequiredPreparationInstruments")
	defer span.End()

	c.logger.Debug("GetAllRequiredPreparationInstruments called")

	return c.querier.GetAllRequiredPreparationInstruments(ctx, results)
}

// GetRequiredPreparationInstruments fetches a list of required preparation instruments from the database that meet a particular filter.
func (c *Client) GetRequiredPreparationInstruments(ctx context.Context, filter *models.QueryFilter) (*models.RequiredPreparationInstrumentList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRequiredPreparationInstruments")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetRequiredPreparationInstruments called")

	requiredPreparationInstrumentList, err := c.querier.GetRequiredPreparationInstruments(ctx, filter)

	return requiredPreparationInstrumentList, err
}

// GetRequiredPreparationInstrumentsWithIDs fetches required preparation instruments from the database within a given set of IDs.
func (c *Client) GetRequiredPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.RequiredPreparationInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRequiredPreparationInstrumentsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetRequiredPreparationInstrumentsWithIDs called")

	requiredPreparationInstrumentList, err := c.querier.GetRequiredPreparationInstrumentsWithIDs(ctx, limit, ids)

	return requiredPreparationInstrumentList, err
}

// CreateRequiredPreparationInstrument creates a required preparation instrument in the database.
func (c *Client) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*models.RequiredPreparationInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRequiredPreparationInstrument")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRequiredPreparationInstrument called")

	return c.querier.CreateRequiredPreparationInstrument(ctx, input)
}

// UpdateRequiredPreparationInstrument updates a particular required preparation instrument. Note that UpdateRequiredPreparationInstrument expects the
// provided input to have a valid ID.
func (c *Client) UpdateRequiredPreparationInstrument(ctx context.Context, updated *models.RequiredPreparationInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRequiredPreparationInstrument")
	defer span.End()

	tracing.AttachRequiredPreparationInstrumentIDToSpan(span, updated.ID)
	c.logger.WithValue("required_preparation_instrument_id", updated.ID).Debug("UpdateRequiredPreparationInstrument called")

	return c.querier.UpdateRequiredPreparationInstrument(ctx, updated)
}

// ArchiveRequiredPreparationInstrument archives a required preparation instrument from the database by its ID.
func (c *Client) ArchiveRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRequiredPreparationInstrument")
	defer span.End()

	tracing.AttachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"required_preparation_instrument_id": requiredPreparationInstrumentID,
	}).Debug("ArchiveRequiredPreparationInstrument called")

	return c.querier.ArchiveRequiredPreparationInstrument(ctx, requiredPreparationInstrumentID)
}

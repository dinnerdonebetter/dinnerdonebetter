package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.ValidInstrumentDataManager = (*Client)(nil)

// ValidInstrumentExists fetches whether or not a valid instrument exists from the database.
func (c *Client) ValidInstrumentExists(ctx context.Context, validInstrumentID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "ValidInstrumentExists")
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"valid_instrument_id": validInstrumentID,
	}).Debug("ValidInstrumentExists called")

	return c.querier.ValidInstrumentExists(ctx, validInstrumentID)
}

// GetValidInstrument fetches a valid instrument from the database.
func (c *Client) GetValidInstrument(ctx context.Context, validInstrumentID uint64) (*models.ValidInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidInstrument")
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"valid_instrument_id": validInstrumentID,
	}).Debug("GetValidInstrument called")

	return c.querier.GetValidInstrument(ctx, validInstrumentID)
}

// GetAllValidInstrumentsCount fetches the count of valid instruments from the database that meet a particular filter.
func (c *Client) GetAllValidInstrumentsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidInstrumentsCount")
	defer span.End()

	c.logger.Debug("GetAllValidInstrumentsCount called")

	return c.querier.GetAllValidInstrumentsCount(ctx)
}

// GetAllValidInstruments fetches a list of all valid instruments in the database.
func (c *Client) GetAllValidInstruments(ctx context.Context, results chan []models.ValidInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidInstruments")
	defer span.End()

	c.logger.Debug("GetAllValidInstruments called")

	return c.querier.GetAllValidInstruments(ctx, results)
}

// GetValidInstruments fetches a list of valid instruments from the database that meet a particular filter.
func (c *Client) GetValidInstruments(ctx context.Context, filter *models.QueryFilter) (*models.ValidInstrumentList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidInstruments")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetValidInstruments called")

	validInstrumentList, err := c.querier.GetValidInstruments(ctx, filter)

	return validInstrumentList, err
}

// GetValidInstrumentsWithIDs fetches valid instruments from the database within a given set of IDs.
func (c *Client) GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidInstrumentsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetValidInstrumentsWithIDs called")

	validInstrumentList, err := c.querier.GetValidInstrumentsWithIDs(ctx, limit, ids)

	return validInstrumentList, err
}

// CreateValidInstrument creates a valid instrument in the database.
func (c *Client) CreateValidInstrument(ctx context.Context, input *models.ValidInstrumentCreationInput) (*models.ValidInstrument, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidInstrument")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateValidInstrument called")

	return c.querier.CreateValidInstrument(ctx, input)
}

// UpdateValidInstrument updates a particular valid instrument. Note that UpdateValidInstrument expects the
// provided input to have a valid ID.
func (c *Client) UpdateValidInstrument(ctx context.Context, updated *models.ValidInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidInstrument")
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, updated.ID)
	c.logger.WithValue("valid_instrument_id", updated.ID).Debug("UpdateValidInstrument called")

	return c.querier.UpdateValidInstrument(ctx, updated)
}

// ArchiveValidInstrument archives a valid instrument from the database by its ID.
func (c *Client) ArchiveValidInstrument(ctx context.Context, validInstrumentID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidInstrument")
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"valid_instrument_id": validInstrumentID,
	}).Debug("ArchiveValidInstrument called")

	return c.querier.ArchiveValidInstrument(ctx, validInstrumentID)
}

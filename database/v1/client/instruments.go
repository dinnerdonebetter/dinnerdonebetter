package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.InstrumentDataManager = (*Client)(nil)

// attachInstrumentIDToSpan provides a consistent way to attach an instrument's ID to a span
func attachInstrumentIDToSpan(span *trace.Span, instrumentID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("instrument_id", strconv.FormatUint(instrumentID, 10)))
	}
}

// GetInstrument fetches an instrument from the database
func (c *Client) GetInstrument(ctx context.Context, instrumentID, userID uint64) (*models.Instrument, error) {
	ctx, span := trace.StartSpan(ctx, "GetInstrument")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachInstrumentIDToSpan(span, instrumentID)

	c.logger.WithValues(map[string]interface{}{
		"instrument_id": instrumentID,
		"user_id":       userID,
	}).Debug("GetInstrument called")

	return c.querier.GetInstrument(ctx, instrumentID, userID)
}

// GetInstrumentCount fetches the count of instruments from the database that meet a particular filter
func (c *Client) GetInstrumentCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetInstrumentCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetInstrumentCount called")

	return c.querier.GetInstrumentCount(ctx, filter, userID)
}

// GetAllInstrumentsCount fetches the count of instruments from the database that meet a particular filter
func (c *Client) GetAllInstrumentsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllInstrumentsCount")
	defer span.End()

	c.logger.Debug("GetAllInstrumentsCount called")

	return c.querier.GetAllInstrumentsCount(ctx)
}

// GetInstruments fetches a list of instruments from the database that meet a particular filter
func (c *Client) GetInstruments(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.InstrumentList, error) {
	ctx, span := trace.StartSpan(ctx, "GetInstruments")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetInstruments called")

	instrumentList, err := c.querier.GetInstruments(ctx, filter, userID)

	return instrumentList, err
}

// GetAllInstrumentsForUser fetches a list of instruments from the database that meet a particular filter
func (c *Client) GetAllInstrumentsForUser(ctx context.Context, userID uint64) ([]models.Instrument, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllInstrumentsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllInstrumentsForUser called")

	instrumentList, err := c.querier.GetAllInstrumentsForUser(ctx, userID)

	return instrumentList, err
}

// CreateInstrument creates an instrument in the database
func (c *Client) CreateInstrument(ctx context.Context, input *models.InstrumentCreationInput) (*models.Instrument, error) {
	ctx, span := trace.StartSpan(ctx, "CreateInstrument")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateInstrument called")

	return c.querier.CreateInstrument(ctx, input)
}

// UpdateInstrument updates a particular instrument. Note that UpdateInstrument expects the
// provided input to have a valid ID.
func (c *Client) UpdateInstrument(ctx context.Context, input *models.Instrument) error {
	ctx, span := trace.StartSpan(ctx, "UpdateInstrument")
	defer span.End()

	attachInstrumentIDToSpan(span, input.ID)
	c.logger.WithValue("instrument_id", input.ID).Debug("UpdateInstrument called")

	return c.querier.UpdateInstrument(ctx, input)
}

// ArchiveInstrument archives an instrument from the database by its ID
func (c *Client) ArchiveInstrument(ctx context.Context, instrumentID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveInstrument")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachInstrumentIDToSpan(span, instrumentID)

	c.logger.WithValues(map[string]interface{}{
		"instrument_id": instrumentID,
		"user_id":       userID,
	}).Debug("ArchiveInstrument called")

	return c.querier.ArchiveInstrument(ctx, instrumentID, userID)
}

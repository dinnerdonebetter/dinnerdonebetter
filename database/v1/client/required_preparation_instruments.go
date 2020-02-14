package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.RequiredPreparationInstrumentDataManager = (*Client)(nil)

// attachRequiredPreparationInstrumentIDToSpan provides a consistent way to attach a required preparation instrument's ID to a span
func attachRequiredPreparationInstrumentIDToSpan(span *trace.Span, requiredPreparationInstrumentID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("required_preparation_instrument_id", strconv.FormatUint(requiredPreparationInstrumentID, 10)))
	}
}

// GetRequiredPreparationInstrument fetches a required preparation instrument from the database
func (c *Client) GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) (*models.RequiredPreparationInstrument, error) {
	ctx, span := trace.StartSpan(ctx, "GetRequiredPreparationInstrument")
	defer span.End()

	attachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"required_preparation_instrument_id": requiredPreparationInstrumentID,
	}).Debug("GetRequiredPreparationInstrument called")

	return c.querier.GetRequiredPreparationInstrument(ctx, requiredPreparationInstrumentID)
}

// GetRequiredPreparationInstrumentCount fetches the count of required preparation instruments from the database that meet a particular filter
func (c *Client) GetRequiredPreparationInstrumentCount(ctx context.Context, filter *models.QueryFilter) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRequiredPreparationInstrumentCount")
	defer span.End()

	attachFilterToSpan(span, filter)

	c.logger.Debug("GetRequiredPreparationInstrumentCount called")

	return c.querier.GetRequiredPreparationInstrumentCount(ctx, filter)
}

// GetAllRequiredPreparationInstrumentsCount fetches the count of required preparation instruments from the database that meet a particular filter
func (c *Client) GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRequiredPreparationInstrumentsCount")
	defer span.End()

	c.logger.Debug("GetAllRequiredPreparationInstrumentsCount called")

	return c.querier.GetAllRequiredPreparationInstrumentsCount(ctx)
}

// GetRequiredPreparationInstruments fetches a list of required preparation instruments from the database that meet a particular filter
func (c *Client) GetRequiredPreparationInstruments(ctx context.Context, filter *models.QueryFilter) (*models.RequiredPreparationInstrumentList, error) {
	ctx, span := trace.StartSpan(ctx, "GetRequiredPreparationInstruments")
	defer span.End()

	attachFilterToSpan(span, filter)

	c.logger.Debug("GetRequiredPreparationInstruments called")

	requiredPreparationInstrumentList, err := c.querier.GetRequiredPreparationInstruments(ctx, filter)

	return requiredPreparationInstrumentList, err
}

// CreateRequiredPreparationInstrument creates a required preparation instrument in the database
func (c *Client) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*models.RequiredPreparationInstrument, error) {
	ctx, span := trace.StartSpan(ctx, "CreateRequiredPreparationInstrument")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRequiredPreparationInstrument called")

	return c.querier.CreateRequiredPreparationInstrument(ctx, input)
}

// UpdateRequiredPreparationInstrument updates a particular required preparation instrument. Note that UpdateRequiredPreparationInstrument expects the
// provided input to have a valid ID.
func (c *Client) UpdateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrument) error {
	ctx, span := trace.StartSpan(ctx, "UpdateRequiredPreparationInstrument")
	defer span.End()

	attachRequiredPreparationInstrumentIDToSpan(span, input.ID)
	c.logger.WithValue("required_preparation_instrument_id", input.ID).Debug("UpdateRequiredPreparationInstrument called")

	return c.querier.UpdateRequiredPreparationInstrument(ctx, input)
}

// ArchiveRequiredPreparationInstrument archives a required preparation instrument from the database by its ID
func (c *Client) ArchiveRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveRequiredPreparationInstrument")
	defer span.End()

	attachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)

	c.logger.WithValues(map[string]interface{}{
		"required_preparation_instrument_id": requiredPreparationInstrumentID,
	}).Debug("ArchiveRequiredPreparationInstrument called")

	return c.querier.ArchiveRequiredPreparationInstrument(ctx, requiredPreparationInstrumentID)
}

package httpclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetValidInstrument gets a valid instrument.
func (c *Client) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	req, err := c.requestBuilder.BuildGetValidInstrumentRequest(ctx, validInstrumentID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get valid instrument request")
	}

	var validInstrument *types.ValidInstrument
	if err = c.fetchAndUnmarshal(ctx, req, &validInstrument); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid instrument")
	}

	return validInstrument, nil
}

// SearchValidInstruments searches through a list of valid instruments.
func (c *Client) SearchValidInstruments(ctx context.Context, query string, limit uint8) ([]*types.ValidInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	if limit == 0 {
		limit = 20
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	req, err := c.requestBuilder.BuildSearchValidInstrumentsRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building search for valid instruments request")
	}

	var validInstruments []*types.ValidInstrument
	if err = c.fetchAndUnmarshal(ctx, req, &validInstruments); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid instruments")
	}

	return validInstruments, nil
}

// GetValidInstruments retrieves a list of valid instruments.
func (c *Client) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (*types.ValidInstrumentList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidInstrumentsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building valid instruments list request")
	}

	var validInstruments *types.ValidInstrumentList
	if err = c.fetchAndUnmarshal(ctx, req, &validInstruments); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid instruments")
	}

	return validInstruments, nil
}

// CreateValidInstrument creates a valid instrument.
func (c *Client) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return "", ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return "", observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidInstrumentRequest(ctx, input)
	if err != nil {
		return "", observability.PrepareError(err, logger, span, "building create valid instrument request")
	}

	var pwr *types.PreWriteResponse
	if err = c.fetchAndUnmarshal(ctx, req, &pwr); err != nil {
		return "", observability.PrepareError(err, logger, span, "creating valid instrument")
	}

	return pwr.ID, nil
}

// UpdateValidInstrument updates a valid instrument.
func (c *Client) UpdateValidInstrument(ctx context.Context, validInstrument *types.ValidInstrument) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validInstrument == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrument.ID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrument.ID)

	req, err := c.requestBuilder.BuildUpdateValidInstrumentRequest(ctx, validInstrument)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update valid instrument request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validInstrument); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid instrument %s", validInstrument.ID)
	}

	return nil
}

// ArchiveValidInstrument archives a valid instrument.
func (c *Client) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	req, err := c.requestBuilder.BuildArchiveValidInstrumentRequest(ctx, validInstrumentID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive valid instrument request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving valid instrument %s", validInstrumentID)
	}

	return nil
}

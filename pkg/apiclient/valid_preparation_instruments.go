package apiclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetValidPreparationInstrument gets a valid ingredient preparation.
func (c *Client) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	req, err := c.requestBuilder.BuildGetValidPreparationInstrumentRequest(ctx, validPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient preparation request")
	}

	var validPreparationInstrument *types.ValidPreparationInstrument
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstrument); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparation")
	}

	return validPreparationInstrument, nil
}

// GetValidPreparationInstruments retrieves a list of valid preparation instruments.
func (c *Client) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (*types.ValidPreparationInstrumentList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidPreparationInstrumentsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation instruments list request")
	}

	var validPreparationInstruments *types.ValidPreparationInstrumentList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation instruments")
	}

	return validPreparationInstruments, nil
}

// GetValidPreparationInstrumentsForPreparation retrieves a list of valid preparation instruments.
func (c *Client) GetValidPreparationInstrumentsForPreparation(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*types.ValidPreparationInstrumentList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	req, err := c.requestBuilder.BuildGetValidPreparationInstrumentsForPreparationRequest(ctx, validPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation instruments list request")
	}

	var validPreparationInstruments *types.ValidPreparationInstrumentList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation instruments")
	}

	return validPreparationInstruments, nil
}

// GetValidPreparationInstrumentsForInstrument retrieves a list of valid preparation instruments.
func (c *Client) GetValidPreparationInstrumentsForInstrument(ctx context.Context, validInstrumentID string, filter *types.QueryFilter) (*types.ValidPreparationInstrumentList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	req, err := c.requestBuilder.BuildGetValidPreparationInstrumentsForInstrumentRequest(ctx, validInstrumentID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation instruments list request")
	}

	var validPreparationInstruments *types.ValidPreparationInstrumentList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation instruments")
	}

	return validPreparationInstruments, nil
}

// CreateValidPreparationInstrument creates a valid ingredient preparation.
func (c *Client) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*types.ValidPreparationInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, input.ID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidPreparationInstrumentRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient preparation request")
	}

	var validPreparationInstrument *types.ValidPreparationInstrument
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstrument); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient preparation")
	}

	return validPreparationInstrument, nil
}

// UpdateValidPreparationInstrument updates a valid ingredient preparation.
func (c *Client) UpdateValidPreparationInstrument(ctx context.Context, validPreparationInstrument *types.ValidPreparationInstrument) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationInstrument == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrument.ID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrument.ID)

	req, err := c.requestBuilder.BuildUpdateValidPreparationInstrumentRequest(ctx, validPreparationInstrument)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstrument); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation %s", validPreparationInstrument.ID)
	}

	return nil
}

// ArchiveValidPreparationInstrument archives a valid ingredient preparation.
func (c *Client) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	req, err := c.requestBuilder.BuildArchiveValidPreparationInstrumentRequest(ctx, validPreparationInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient preparation %s", validPreparationInstrumentID)
	}

	return nil
}

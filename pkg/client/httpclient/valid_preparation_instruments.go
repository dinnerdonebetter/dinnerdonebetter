package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// ValidPreparationInstrumentExists retrieves whether a valid preparation instrument exists.
func (c *Client) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationInstrumentID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	req, err := c.requestBuilder.BuildValidPreparationInstrumentExistsRequest(ctx, validPreparationInstrumentID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building valid preparation instrument existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for valid preparation instrument #%d", validPreparationInstrumentID)
	}

	return exists, nil
}

// GetValidPreparationInstrument gets a valid preparation instrument.
func (c *Client) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) (*types.ValidPreparationInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	req, err := c.requestBuilder.BuildGetValidPreparationInstrumentRequest(ctx, validPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get valid preparation instrument request")
	}

	var validPreparationInstrument *types.ValidPreparationInstrument
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstrument); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid preparation instrument")
	}

	return validPreparationInstrument, nil
}

// GetValidPreparationInstruments retrieves a list of valid preparation instruments.
func (c *Client) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (*types.ValidPreparationInstrumentList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidPreparationInstrumentsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building valid preparation instruments list request")
	}

	var validPreparationInstruments *types.ValidPreparationInstrumentList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid preparation instruments")
	}

	return validPreparationInstruments, nil
}

// CreateValidPreparationInstrument creates a valid preparation instrument.
func (c *Client) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationInput) (*types.ValidPreparationInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidPreparationInstrumentRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create valid preparation instrument request")
	}

	var validPreparationInstrument *types.ValidPreparationInstrument
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstrument); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating valid preparation instrument")
	}

	return validPreparationInstrument, nil
}

// UpdateValidPreparationInstrument updates a valid preparation instrument.
func (c *Client) UpdateValidPreparationInstrument(ctx context.Context, validPreparationInstrument *types.ValidPreparationInstrument) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationInstrument == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrument.ID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrument.ID)

	req, err := c.requestBuilder.BuildUpdateValidPreparationInstrumentRequest(ctx, validPreparationInstrument)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update valid preparation instrument request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstrument); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid preparation instrument #%d", validPreparationInstrument.ID)
	}

	return nil
}

// ArchiveValidPreparationInstrument archives a valid preparation instrument.
func (c *Client) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationInstrumentID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	req, err := c.requestBuilder.BuildArchiveValidPreparationInstrumentRequest(ctx, validPreparationInstrumentID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive valid preparation instrument request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving valid preparation instrument #%d", validPreparationInstrumentID)
	}

	return nil
}

// GetAuditLogForValidPreparationInstrument retrieves a list of audit log entries pertaining to a valid preparation instrument.
func (c *Client) GetAuditLogForValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	req, err := c.requestBuilder.BuildGetAuditLogForValidPreparationInstrumentRequest(ctx, validPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for valid preparation instrument request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

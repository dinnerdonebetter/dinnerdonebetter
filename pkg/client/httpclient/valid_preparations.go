package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// ValidPreparationExists retrieves whether a valid preparation exists.
func (c *Client) ValidPreparationExists(ctx context.Context, validPreparationID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	req, err := c.requestBuilder.BuildValidPreparationExistsRequest(ctx, validPreparationID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building valid preparation existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for valid preparation #%d", validPreparationID)
	}

	return exists, nil
}

// GetValidPreparation gets a valid preparation.
func (c *Client) GetValidPreparation(ctx context.Context, validPreparationID uint64) (*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	req, err := c.requestBuilder.BuildGetValidPreparationRequest(ctx, validPreparationID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get valid preparation request")
	}

	var validPreparation *types.ValidPreparation
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparation); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid preparation")
	}

	return validPreparation, nil
}

// SearchValidPreparations searches through a list of valid preparations.
func (c *Client) SearchValidPreparations(ctx context.Context, query string, limit uint8) ([]*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	if limit == 0 {
		limit = 20
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	req, err := c.requestBuilder.BuildSearchValidPreparationsRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building search for valid preparations request")
	}

	var validPreparations []*types.ValidPreparation
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparations); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid preparations")
	}

	return validPreparations, nil
}

// GetValidPreparations retrieves a list of valid preparations.
func (c *Client) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (*types.ValidPreparationList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidPreparationsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building valid preparations list request")
	}

	var validPreparations *types.ValidPreparationList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparations); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid preparations")
	}

	return validPreparations, nil
}

// CreateValidPreparation creates a valid preparation.
func (c *Client) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationInput) (*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidPreparationRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create valid preparation request")
	}

	var validPreparation *types.ValidPreparation
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparation); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating valid preparation")
	}

	return validPreparation, nil
}

// UpdateValidPreparation updates a valid preparation.
func (c *Client) UpdateValidPreparation(ctx context.Context, validPreparation *types.ValidPreparation) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparation == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparation.ID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparation.ID)

	req, err := c.requestBuilder.BuildUpdateValidPreparationRequest(ctx, validPreparation)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update valid preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validPreparation); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid preparation #%d", validPreparation.ID)
	}

	return nil
}

// ArchiveValidPreparation archives a valid preparation.
func (c *Client) ArchiveValidPreparation(ctx context.Context, validPreparationID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	req, err := c.requestBuilder.BuildArchiveValidPreparationRequest(ctx, validPreparationID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive valid preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving valid preparation #%d", validPreparationID)
	}

	return nil
}

// GetAuditLogForValidPreparation retrieves a list of audit log entries pertaining to a valid preparation.
func (c *Client) GetAuditLogForValidPreparation(ctx context.Context, validPreparationID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	req, err := c.requestBuilder.BuildGetAuditLogForValidPreparationRequest(ctx, validPreparationID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for valid preparation request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

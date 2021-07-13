package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// ValidIngredientPreparationExists retrieves whether a valid ingredient preparation exists.
func (c *Client) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientPreparationID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	req, err := c.requestBuilder.BuildValidIngredientPreparationExistsRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building valid ingredient preparation existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for valid ingredient preparation #%d", validIngredientPreparationID)
	}

	return exists, nil
}

// GetValidIngredientPreparation gets a valid ingredient preparation.
func (c *Client) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) (*types.ValidIngredientPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	req, err := c.requestBuilder.BuildGetValidIngredientPreparationRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get valid ingredient preparation request")
	}

	var validIngredientPreparation *types.ValidIngredientPreparation
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientPreparation); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid ingredient preparation")
	}

	return validIngredientPreparation, nil
}

// GetValidIngredientPreparations retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientPreparations(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientPreparationList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidIngredientPreparationsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building valid ingredient preparations list request")
	}

	var validIngredientPreparations *types.ValidIngredientPreparationList
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientPreparations); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid ingredient preparations")
	}

	return validIngredientPreparations, nil
}

// CreateValidIngredientPreparation creates a valid ingredient preparation.
func (c *Client) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationInput) (*types.ValidIngredientPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidIngredientPreparationRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create valid ingredient preparation request")
	}

	var validIngredientPreparation *types.ValidIngredientPreparation
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientPreparation); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating valid ingredient preparation")
	}

	return validIngredientPreparation, nil
}

// UpdateValidIngredientPreparation updates a valid ingredient preparation.
func (c *Client) UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparation *types.ValidIngredientPreparation) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientPreparation == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparation.ID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparation.ID)

	req, err := c.requestBuilder.BuildUpdateValidIngredientPreparationRequest(ctx, validIngredientPreparation)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update valid ingredient preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientPreparation); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient preparation #%d", validIngredientPreparation.ID)
	}

	return nil
}

// ArchiveValidIngredientPreparation archives a valid ingredient preparation.
func (c *Client) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientPreparationID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	req, err := c.requestBuilder.BuildArchiveValidIngredientPreparationRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive valid ingredient preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving valid ingredient preparation #%d", validIngredientPreparationID)
	}

	return nil
}

// GetAuditLogForValidIngredientPreparation retrieves a list of audit log entries pertaining to a valid ingredient preparation.
func (c *Client) GetAuditLogForValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	req, err := c.requestBuilder.BuildGetAuditLogForValidIngredientPreparationRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for valid ingredient preparation request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

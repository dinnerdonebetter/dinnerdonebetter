package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

// GetValidIngredientPreparation gets a valid ingredient preparation.
func (c *Client) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	req, err := c.requestBuilder.BuildGetValidIngredientPreparationRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient preparation request")
	}

	var validIngredientPreparation *types.ValidIngredientPreparation
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientPreparation); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparation")
	}

	return validIngredientPreparation, nil
}

// GetValidIngredientPreparations retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientPreparations(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientPreparationList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidIngredientPreparationsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient preparations list request")
	}

	var validIngredientPreparations *types.ValidIngredientPreparationList
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientPreparations); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparations")
	}

	return validIngredientPreparations, nil
}

// GetValidIngredientPreparationsForIngredient retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientPreparationsForIngredient(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.ValidIngredientPreparationList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildGetValidIngredientPreparationsForIngredientRequest(ctx, validIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient preparations list request")
	}

	var validPreparationInstruments *types.ValidIngredientPreparationList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparations")
	}

	return validPreparationInstruments, nil
}

// GetValidIngredientPreparationsForPreparation retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientPreparationsForPreparation(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*types.ValidIngredientPreparationList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	req, err := c.requestBuilder.BuildGetValidIngredientPreparationsForPreparationRequest(ctx, validPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient preparations list request")
	}

	var validPreparationInstruments *types.ValidIngredientPreparationList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparations")
	}

	return validPreparationInstruments, nil
}

// CreateValidIngredientPreparation creates a valid ingredient preparation.
func (c *Client) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationRequestInput) (*types.ValidIngredientPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidIngredientPreparationRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient preparation request")
	}

	var validIngredientPreparation *types.ValidIngredientPreparation
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientPreparation); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient preparation")
	}

	return validIngredientPreparation, nil
}

// UpdateValidIngredientPreparation updates a valid ingredient preparation.
func (c *Client) UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparation *types.ValidIngredientPreparation) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientPreparation == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparation.ID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparation.ID)

	req, err := c.requestBuilder.BuildUpdateValidIngredientPreparationRequest(ctx, validIngredientPreparation)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientPreparation); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation %s", validIngredientPreparation.ID)
	}

	return nil
}

// ArchiveValidIngredientPreparation archives a valid ingredient preparation.
func (c *Client) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	req, err := c.requestBuilder.BuildArchiveValidIngredientPreparationRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient preparation %s", validIngredientPreparationID)
	}

	return nil
}

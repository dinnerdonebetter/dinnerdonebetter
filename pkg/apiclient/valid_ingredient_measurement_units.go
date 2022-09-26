package apiclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetValidIngredientMeasurementUnit gets a valid ingredient preparation.
func (c *Client) GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, validIngredientMeasurementUnitID)

	req, err := c.requestBuilder.BuildGetValidIngredientMeasurementUnitRequest(ctx, validIngredientMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient preparation request")
	}

	var validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientMeasurementUnit); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparation")
	}

	return validIngredientMeasurementUnit, nil
}

// GetValidIngredientMeasurementUnits retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientMeasurementUnitList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidIngredientMeasurementUnitsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient preparations list request")
	}

	var validIngredientMeasurementUnits *types.ValidIngredientMeasurementUnitList
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientMeasurementUnits); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparations")
	}

	return validIngredientMeasurementUnits, nil
}

// GetValidIngredientMeasurementUnitsForIngredient retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientMeasurementUnitsForIngredient(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.ValidIngredientMeasurementUnitList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildGetValidIngredientMeasurementUnitsForIngredientRequest(ctx, validIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient preparations list request")
	}

	var validPreparationInstruments *types.ValidIngredientMeasurementUnitList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparations")
	}

	return validPreparationInstruments, nil
}

// GetValidIngredientMeasurementUnitsForMeasurementUnit retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx context.Context, validMeasurementUnitID string, filter *types.QueryFilter) (*types.ValidIngredientMeasurementUnitList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	req, err := c.requestBuilder.BuildGetValidIngredientMeasurementUnitsForMeasurementUnitRequest(ctx, validMeasurementUnitID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient preparations list request")
	}

	var validPreparationInstruments *types.ValidIngredientMeasurementUnitList
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparations")
	}

	return validPreparationInstruments, nil
}

// CreateValidIngredientMeasurementUnit creates a valid ingredient preparation.
func (c *Client) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitCreationRequestInput) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidIngredientMeasurementUnitRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient preparation request")
	}

	var validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientMeasurementUnit); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient preparation")
	}

	return validIngredientMeasurementUnit, nil
}

// UpdateValidIngredientMeasurementUnit updates a valid ingredient preparation.
func (c *Client) UpdateValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientMeasurementUnit == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnit.ID)
	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, validIngredientMeasurementUnit.ID)

	req, err := c.requestBuilder.BuildUpdateValidIngredientMeasurementUnitRequest(ctx, validIngredientMeasurementUnit)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientMeasurementUnit); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation %s", validIngredientMeasurementUnit.ID)
	}

	return nil
}

// ArchiveValidIngredientMeasurementUnit archives a valid ingredient preparation.
func (c *Client) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, validIngredientMeasurementUnitID)

	req, err := c.requestBuilder.BuildArchiveValidIngredientMeasurementUnitRequest(ctx, validIngredientMeasurementUnitID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient preparation %s", validIngredientMeasurementUnitID)
	}

	return nil
}

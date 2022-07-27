package httpclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetValidMeasurementUnit gets a valid ingredient.
func (c *Client) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	req, err := c.requestBuilder.BuildGetValidMeasurementUnitRequest(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get valid ingredient request")
	}

	var validMeasurementUnit *types.ValidMeasurementUnit
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnit); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid ingredient")
	}

	return validMeasurementUnit, nil
}

// SearchValidMeasurementUnits searches through a list of valid ingredients.
func (c *Client) SearchValidMeasurementUnits(ctx context.Context, query string, limit uint8) ([]*types.ValidMeasurementUnit, error) {
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

	req, err := c.requestBuilder.BuildSearchValidMeasurementUnitsRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building search for valid ingredients request")
	}

	var validMeasurementUnits []*types.ValidMeasurementUnit
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnits); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid ingredients")
	}

	return validMeasurementUnits, nil
}

// GetValidMeasurementUnits retrieves a list of valid ingredients.
func (c *Client) GetValidMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (*types.ValidMeasurementUnitList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidMeasurementUnitsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building valid ingredients list request")
	}

	var validMeasurementUnits *types.ValidMeasurementUnitList
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnits); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid ingredients")
	}

	return validMeasurementUnits, nil
}

// CreateValidMeasurementUnit creates a valid ingredient.
func (c *Client) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitCreationRequestInput) (*types.ValidMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidMeasurementUnitRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create valid ingredient request")
	}

	var validMeasurementUnit *types.ValidMeasurementUnit
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnit); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating valid ingredient")
	}

	return validMeasurementUnit, nil
}

// UpdateValidMeasurementUnit updates a valid ingredient.
func (c *Client) UpdateValidMeasurementUnit(ctx context.Context, validMeasurementUnit *types.ValidMeasurementUnit) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnit == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnit.ID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnit.ID)

	req, err := c.requestBuilder.BuildUpdateValidMeasurementUnitRequest(ctx, validMeasurementUnit)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update valid ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnit); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient %s", validMeasurementUnit.ID)
	}

	return nil
}

// ArchiveValidMeasurementUnit archives a valid ingredient.
func (c *Client) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	req, err := c.requestBuilder.BuildArchiveValidMeasurementUnitRequest(ctx, validMeasurementUnitID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive valid ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving valid ingredient %s", validMeasurementUnitID)
	}

	return nil
}

package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidMeasurementUnitConversion gets a valid measurement conversion.
func (c *Client) GetValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	req, err := c.requestBuilder.BuildGetValidMeasurementUnitConversionRequest(ctx, validMeasurementUnitConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid measurement conversion request")
	}

	var validMeasurementUnitConversion *types.ValidMeasurementUnitConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnitConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	return validMeasurementUnitConversion, nil
}

// GetValidMeasurementUnitConversionsFromUnit gets a valid measurement conversion.
func (c *Client) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	req, err := c.requestBuilder.BuildGetValidMeasurementUnitConversionsFromUnitRequest(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid measurement conversion request")
	}

	var validMeasurementUnitConversion []*types.ValidMeasurementUnitConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnitConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	return validMeasurementUnitConversion, nil
}

// GetValidMeasurementUnitConversionToUnit gets a valid measurement conversion.
func (c *Client) GetValidMeasurementUnitConversionToUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	req, err := c.requestBuilder.BuildGetValidMeasurementUnitConversionsToUnitRequest(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid measurement conversion request")
	}

	var validMeasurementUnitConversion []*types.ValidMeasurementUnitConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnitConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	return validMeasurementUnitConversion, nil
}

// CreateValidMeasurementUnitConversion creates a valid measurement conversion.
func (c *Client) CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionCreationRequestInput) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidMeasurementUnitConversionRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid measurement conversion request")
	}

	var validMeasurementUnitConversion *types.ValidMeasurementUnitConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnitConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid measurement conversion")
	}

	return validMeasurementUnitConversion, nil
}

// UpdateValidMeasurementUnitConversion updates a valid measurement conversion.
func (c *Client) UpdateValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversion *types.ValidMeasurementUnitConversion) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitConversion == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversion.ID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversion.ID)

	req, err := c.requestBuilder.BuildUpdateValidMeasurementUnitConversionRequest(ctx, validMeasurementUnitConversion)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid measurement conversion request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementUnitConversion); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement conversion %s", validMeasurementUnitConversion.ID)
	}

	return nil
}

// ArchiveValidMeasurementUnitConversion archives a valid measurement conversion.
func (c *Client) ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	req, err := c.requestBuilder.BuildArchiveValidMeasurementUnitConversionRequest(ctx, validMeasurementUnitConversionID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid measurement conversion request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement conversion %s", validMeasurementUnitConversionID)
	}

	return nil
}

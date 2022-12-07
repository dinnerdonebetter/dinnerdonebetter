package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

// GetValidMeasurementConversion gets a valid measurement conversion.
func (c *Client) GetValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) (*types.ValidMeasurementConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)

	req, err := c.requestBuilder.BuildGetValidMeasurementConversionRequest(ctx, validMeasurementConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid measurement conversion request")
	}

	var validMeasurementConversion *types.ValidMeasurementConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	return validMeasurementConversion, nil
}

// GetValidMeasurementConversionsFromUnit gets a valid measurement conversion.
func (c *Client) GetValidMeasurementConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	req, err := c.requestBuilder.BuildGetValidMeasurementConversionsFromUnitRequest(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid measurement conversion request")
	}

	var validMeasurementConversion []*types.ValidMeasurementConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	return validMeasurementConversion, nil
}

// GetValidMeasurementConversionToUnit gets a valid measurement conversion.
func (c *Client) GetValidMeasurementConversionToUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	req, err := c.requestBuilder.BuildGetValidMeasurementConversionsToUnitRequest(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid measurement conversion request")
	}

	var validMeasurementConversion []*types.ValidMeasurementConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	return validMeasurementConversion, nil
}

// CreateValidMeasurementConversion creates a valid measurement conversion.
func (c *Client) CreateValidMeasurementConversion(ctx context.Context, input *types.ValidMeasurementConversionCreationRequestInput) (*types.ValidMeasurementConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidMeasurementConversionRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid measurement conversion request")
	}

	var validMeasurementConversion *types.ValidMeasurementConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid measurement conversion")
	}

	return validMeasurementConversion, nil
}

// UpdateValidMeasurementConversion updates a valid measurement conversion.
func (c *Client) UpdateValidMeasurementConversion(ctx context.Context, validMeasurementConversion *types.ValidMeasurementConversion) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementConversion == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversion.ID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversion.ID)

	req, err := c.requestBuilder.BuildUpdateValidMeasurementConversionRequest(ctx, validMeasurementConversion)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid measurement conversion request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementConversion); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement conversion %s", validMeasurementConversion.ID)
	}

	return nil
}

// ArchiveValidMeasurementConversion archives a valid measurement conversion.
func (c *Client) ArchiveValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementConversionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)

	req, err := c.requestBuilder.BuildArchiveValidMeasurementConversionRequest(ctx, validMeasurementConversionID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid measurement conversion request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement conversion %s", validMeasurementConversionID)
	}

	return nil
}

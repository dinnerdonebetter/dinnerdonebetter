package apiclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetValidMeasurementConversion gets a valid preparation.
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
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid preparation request")
	}

	var validMeasurementConversion *types.ValidMeasurementConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation")
	}

	return validMeasurementConversion, nil
}

// CreateValidMeasurementConversion creates a valid preparation.
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
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid preparation request")
	}

	var validMeasurementConversion *types.ValidMeasurementConversion
	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementConversion); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid preparation")
	}

	return validMeasurementConversion, nil
}

// UpdateValidMeasurementConversion updates a valid preparation.
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
		return observability.PrepareAndLogError(err, logger, span, "building update valid preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validMeasurementConversion); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation %s", validMeasurementConversion.ID)
	}

	return nil
}

// ArchiveValidMeasurementConversion archives a valid preparation.
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
		return observability.PrepareAndLogError(err, logger, span, "building archive valid preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid preparation %s", validMeasurementConversionID)
	}

	return nil
}

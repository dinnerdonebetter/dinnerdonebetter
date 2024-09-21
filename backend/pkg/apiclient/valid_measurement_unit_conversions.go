package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
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
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)

	res, err := c.authedGeneratedClient.GetValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get valid measurement conversion")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnitConversion]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	res, err := c.authedGeneratedClient.GetValidMeasurementUnitConversionsFromUnit(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get valid measurement conversion")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidMeasurementUnitConversion]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	res, err := c.authedGeneratedClient.ValidMeasurementUnitConversionsToUnit(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get valid measurement conversion")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidMeasurementUnitConversion]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement conversion")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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

	body := generated.CreateValidMeasurementUnitConversionJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidMeasurementUnitConversion(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create valid measurement conversion")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnitConversion]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid measurement conversion")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversion.ID)

	body := generated.UpdateValidMeasurementUnitConversionJSONRequestBody{}
	c.copyType(&body, validMeasurementUnitConversion)

	res, err := c.authedGeneratedClient.UpdateValidMeasurementUnitConversion(ctx, validMeasurementUnitConversion.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update valid measurement conversion")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnitConversion]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement conversion %s", validMeasurementUnitConversion.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
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
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)

	res, err := c.authedGeneratedClient.ArchiveValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive valid measurement conversion")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnitConversion]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement conversion %s", validMeasurementUnitConversionID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) UpdateValidMeasurementUnitConversion(
	ctx context.Context,
	validMeasurementUnitConversionID string,
	input *types.ValidMeasurementUnitConversionUpdateRequestInput,
	reqMods ...RequestModifier,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_measurement_conversions/%s", validMeasurementUnitConversionID))
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a ValidMeasurementUnitConversion")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnitConversion]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading ValidMeasurementUnitConversion creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

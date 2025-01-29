// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
)

func (c *Client) CreateValidMeasurementUnitConversion(
	ctx context.Context,
	input *ValidMeasurementUnitConversionCreationRequestInput,
	reqMods ...RequestModifier,
) (*ValidMeasurementUnitConversion, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	u := c.BuildURL(ctx, nil, "/api/v1/valid_measurement_conversions")
	req, err := c.buildDataRequest(ctx, http.MethodPost, u, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a ValidMeasurementUnitConversion")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*ValidMeasurementUnitConversion]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading ValidMeasurementUnitConversion creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

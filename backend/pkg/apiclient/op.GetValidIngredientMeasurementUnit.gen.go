// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

func (c *Client) GetValidIngredientMeasurementUnit(
	ctx context.Context,
	validIngredientMeasurementUnitID string,
	reqMods ...RequestModifier,
) (*ValidIngredientMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return nil, buildInvalidIDError("validIngredientMeasurementUnit")
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_ingredient_measurement_units/%s", validIngredientMeasurementUnitID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a ValidIngredientMeasurementUnit")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*ValidIngredientMeasurementUnit]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading ValidIngredientMeasurementUnit response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

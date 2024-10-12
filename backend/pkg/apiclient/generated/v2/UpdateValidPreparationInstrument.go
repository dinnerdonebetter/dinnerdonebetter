// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient




import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
)


func (c *Client) UpdateValidPreparationInstrument(
	ctx context.Context,
validPreparationVesselID string,
input *types.ValidPreparationInstrumentUpdateRequestInput,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVesselID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, validPreparationVesselID)

 


	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_preparation_instruments/%s" , validPreparationVesselID ))
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a ValidPreparationInstrument")
	}

	var apiResponse *types.APIResponse[ *types.ValidPreparationInstrument]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading ValidPreparationInstrument creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}


	return nil
}
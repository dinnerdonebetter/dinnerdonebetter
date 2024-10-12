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


func (c *Client) GetValidVessel(
	ctx context.Context,
validVesselID string,
) ( *types.ValidVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validVesselID == "" {
		return nil, buildInvalidIDError("validVessel")
	} 
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_vessels/%s" , validVesselID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a ValidVessel")
	}

	var apiResponse *types.APIResponse[  *types.ValidVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading ValidVessel response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}


	return apiResponse.Data, nil
}
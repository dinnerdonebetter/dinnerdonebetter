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


func (c *Client) UpdateValidPreparation(
	ctx context.Context,
validPreparationID string,
input *types.ValidPreparationUpdateRequestInput,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

 


	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_preparations/%s" , validPreparationID ))
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a ValidPreparation")
	}

	var apiResponse *types.APIResponse[ *types.ValidPreparation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading ValidPreparation creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}


	return nil
}
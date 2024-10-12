// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient




import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
)


func (c *Client) UpdateUserDetails(
	ctx context.Context,
input *types.UserDetailsUpdateRequestInput,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

 


	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/users/details" ))
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a User")
	}

	var apiResponse *types.APIResponse[ *types.User]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading User creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}


	return nil
}
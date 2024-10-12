// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient




import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
)


func (c *Client) RunMealPlanTaskCreatorWorker(
	ctx context.Context,
input *types.CreateMealPlanTasksRequest,
) (*types.CreateMealPlanTasksResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}


	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/workers/meal_plan_tasks" ))
	req, err := c.buildDataRequest(ctx, http.MethodPost, u, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a CreateMealPlanTasksResponse")
	}

	var apiResponse *types.APIResponse[ *types.CreateMealPlanTasksResponse]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading CreateMealPlanTasksResponse creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}


	return apiResponse.Data, nil
}
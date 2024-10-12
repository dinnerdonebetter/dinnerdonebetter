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


func (c *Client) UpdateMealPlanOptionVote(
	ctx context.Context,
mealPlanID string,
mealPlanEventID string,
mealPlanOptionID string,
mealPlanOptionVoteID string,
input *types.MealPlanOptionVoteUpdateRequestInput,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	if mealPlanOptionVoteID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

 


	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s" , mealPlanID , mealPlanEventID , mealPlanOptionID , mealPlanOptionVoteID ))
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a MealPlanOptionVote")
	}

	var apiResponse *types.APIResponse[ *types.MealPlanOptionVote]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading MealPlanOptionVote creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}


	return nil
}
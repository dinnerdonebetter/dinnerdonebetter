// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) UpdateMealPlanTaskStatus(
	ctx context.Context,
	mealPlanID string,
	mealPlanTaskID string,
	input *types.MealPlanTaskStatusChangeRequestInput,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanTaskID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, mealPlanTaskID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/meal_plans/%s/tasks/%s", mealPlanID, mealPlanTaskID))
	req, err := c.buildDataRequest(ctx, http.MethodPatch, u, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a MealPlanTask")
	}

	var apiResponse *types.APIResponse[*types.MealPlanTask]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading MealPlanTask creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

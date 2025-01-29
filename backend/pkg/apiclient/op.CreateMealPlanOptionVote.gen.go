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

func (c *Client) CreateMealPlanOptionVote(
	ctx context.Context,
	mealPlanID string,
	mealPlanEventID string,
	input *MealPlanOptionVoteCreationRequestInput,
	reqMods ...RequestModifier,
) ([]*MealPlanOptionVote, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if mealPlanID == "" {
		return nil, buildInvalidIDError("mealPlan")
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, buildInvalidIDError("mealPlanEvent")
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/meal_plans/%s/events/%s/vote", mealPlanID, mealPlanEventID))
	req, err := c.buildDataRequest(ctx, http.MethodPost, u, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to create a MealPlanOptionVote")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[[]*MealPlanOptionVote]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading MealPlanOptionVote creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

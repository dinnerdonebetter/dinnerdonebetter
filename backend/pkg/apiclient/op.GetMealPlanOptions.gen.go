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

func (c *Client) GetMealPlanOptions(
	ctx context.Context,
	mealPlanID string,
	mealPlanEventID string,
	filter *QueryFilter,
	reqMods ...RequestModifier,
) (*QueryFilteredResult[MealPlanOption], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

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

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, fmt.Sprintf("/api/v1/meal_plans/%s/events/%s/options", mealPlanID, mealPlanEventID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of MealPlanOption")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[[]*MealPlanOption]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of MealPlanOption")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &QueryFilteredResult[MealPlanOption]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

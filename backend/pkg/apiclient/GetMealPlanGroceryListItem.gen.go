// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) GetMealPlanGroceryListItem(
	ctx context.Context,
	mealPlanID string,
	mealPlanGroceryListItemID string,
) (*types.MealPlanGroceryListItem, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, buildInvalidIDError("mealPlan")
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return nil, buildInvalidIDError("mealPlanGroceryListItem")
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/meal_plans/%s/grocery_list_items/%s", mealPlanID, mealPlanGroceryListItemID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a MealPlanGroceryListItem")
	}

	var apiResponse *types.APIResponse[*types.MealPlanGroceryListItem]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading MealPlanGroceryListItem response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}
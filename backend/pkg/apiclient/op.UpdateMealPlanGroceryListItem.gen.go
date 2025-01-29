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

func (c *Client) UpdateMealPlanGroceryListItem(
	ctx context.Context,
	mealPlanID string,
	mealPlanGroceryListItemID string,
	input *MealPlanGroceryListItemUpdateRequestInput,
	reqMods ...RequestModifier,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/meal_plans/%s/grocery_list_items/%s", mealPlanID, mealPlanGroceryListItemID))
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a MealPlanGroceryListItem")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*MealPlanGroceryListItem]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading MealPlanGroceryListItem creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetMealPlanGroceryListItem gets a meal plan grocery list item.
func (c *Client) GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItemID)

	req, err := c.requestBuilder.BuildGetMealPlanGroceryListItemRequest(ctx, mealPlanID, mealPlanGroceryListItemID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal plan grocery list item request")
	}

	var validIngredient *types.MealPlanGroceryListItem
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan grocery list item")
	}

	return validIngredient, nil
}

// GetMealPlanGroceryListItemsForMealPlan gets a meal plan grocery list item.
func (c *Client) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) ([]*types.MealPlanGroceryListItem, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	req, err := c.requestBuilder.BuildGetMealPlanGroceryListItemsForMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal plan grocery list item request")
	}

	var mealPlanGroceryListItems []*types.MealPlanGroceryListItem
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanGroceryListItems); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan grocery list item")
	}

	return mealPlanGroceryListItems, nil
}

// CreateMealPlanGroceryListItem creates a meal plan grocery list item.
func (c *Client) CreateMealPlanGroceryListItem(ctx context.Context, mealPlanID string, input *types.MealPlanGroceryListItemCreationRequestInput) (*types.MealPlanGroceryListItem, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if input == nil {
		return nil, ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildCreateMealPlanGroceryListItemRequest(ctx, mealPlanID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal plan grocery list item request")
	}

	var validIngredient *types.MealPlanGroceryListItem
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan grocery list item")
	}

	return validIngredient, nil
}

// UpdateMealPlanGroceryListItem updates a meal plan grocery list item.
func (c *Client) UpdateMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string, input *types.MealPlanGroceryListItemUpdateRequestInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return ErrNilInputProvided
	}

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildUpdateMealPlanGroceryListItemRequest(ctx, mealPlanID, mealPlanGroceryListItemID, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building create meal plan grocery list item request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "creating meal plan grocery list item")
	}

	return nil
}

// ArchiveMealPlanGroceryListItem updates a meal plan grocery list item.
func (c *Client) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItemID)

	req, err := c.requestBuilder.BuildArchiveMealPlanGroceryListItemRequest(ctx, mealPlanID, mealPlanGroceryListItemID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building create meal plan grocery list item request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan %s", mealPlanID)
	}

	return nil
}

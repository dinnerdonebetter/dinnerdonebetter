package httpclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetMeal gets a meal.
func (c *Client) GetMeal(ctx context.Context, mealID string) (*types.Meal, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	req, err := c.requestBuilder.BuildGetMealRequest(ctx, mealID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get meal request")
	}

	var meal *types.Meal
	if err = c.fetchAndUnmarshal(ctx, req, &meal); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving meal")
	}

	return meal, nil
}

// GetMeals retrieves a list of meals.
func (c *Client) GetMeals(ctx context.Context, filter *types.QueryFilter) (*types.MealList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetMealsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building meals list request")
	}

	var meals *types.MealList
	if err = c.fetchAndUnmarshal(ctx, req, &meals); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving meals")
	}

	return meals, nil
}

// CreateMeal creates a meal.
func (c *Client) CreateMeal(ctx context.Context, input *types.MealCreationRequestInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return "", ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return "", observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealRequest(ctx, input)
	if err != nil {
		return "", observability.PrepareError(err, logger, span, "building create meal request")
	}

	var pwr *types.PreWriteResponse
	if err = c.fetchAndUnmarshal(ctx, req, &pwr); err != nil {
		return "", observability.PrepareError(err, logger, span, "creating meal")
	}

	return pwr.ID, nil
}

// ArchiveMeal archives a meal.
func (c *Client) ArchiveMeal(ctx context.Context, mealID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	req, err := c.requestBuilder.BuildArchiveMealRequest(ctx, mealID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive meal request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving meal %s", mealID)
	}

	return nil
}

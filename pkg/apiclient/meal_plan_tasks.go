package apiclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetMealPlanTask gets a meal plan task.
func (c *Client) GetMealPlanTask(ctx context.Context, mealPlanID, mealPlanTaskID string) (*types.MealPlanTask, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachMealPlanTaskIDToSpan(span, mealPlanTaskID)

	req, err := c.requestBuilder.BuildGetMealPlanTaskRequest(ctx, mealPlanID, mealPlanTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal plan task request")
	}

	var validIngredient *types.MealPlanTask
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan task")
	}

	return validIngredient, nil
}

// CreateMealPlanTask creates a meal plan task.
func (c *Client) CreateMealPlanTask(ctx context.Context, mealPlanID string, input *types.MealPlanTaskCreationRequestInput) (*types.MealPlanTask, error) {
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

	req, err := c.requestBuilder.BuildCreateMealPlanTaskRequest(ctx, mealPlanID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal plan task request")
	}

	var validIngredient *types.MealPlanTask
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan task")
	}

	return validIngredient, nil
}

// UpdateMealPlanTaskStatus updates a meal plan task.
func (c *Client) UpdateMealPlanTaskStatus(ctx context.Context, mealPlanID string, input *types.MealPlanTaskStatusChangeRequestInput) (*types.MealPlanTask, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildChangeMealPlanTaskStatusRequest(ctx, mealPlanID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create meal plan task request")
	}

	var validIngredient *types.MealPlanTask
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	return validIngredient, nil
}

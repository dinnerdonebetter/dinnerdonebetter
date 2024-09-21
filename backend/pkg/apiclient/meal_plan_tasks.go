package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
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
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, mealPlanTaskID)

	res, err := c.authedGeneratedClient.GetMealPlanTask(ctx, mealPlanID, mealPlanTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get meal plan task")
	}

	var apiResponse *types.APIResponse[*types.MealPlanTask]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan task")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if input == nil {
		return nil, ErrInvalidIDProvided
	}

	body := generated.CreateMealPlanTaskJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateMealPlanTask(ctx, mealPlanID, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get meal plan task")
	}

	var apiResponse *types.APIResponse[*types.MealPlanTask]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan task")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.UpdateMealPlanTaskStatusJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.UpdateMealPlanTaskStatus(ctx, mealPlanID, input.ID, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create meal plan task")
	}

	var apiResponse *types.APIResponse[*types.MealPlanTask]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

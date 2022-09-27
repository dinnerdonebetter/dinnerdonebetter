package apiclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetMealPlanTask gets an advanced prep step.
func (c *Client) GetMealPlanTask(ctx context.Context, validIngredientID string) (*types.MealPlanTask, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, validIngredientID)
	tracing.AttachMealPlanTaskIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildGetMealPlanTaskRequest(ctx, validIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get advanced prep step request")
	}

	var validIngredient *types.MealPlanTask
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving advanced prep step")
	}

	return validIngredient, nil
}

// UpdateMealPlanTaskStatus updates an advanced prep step.
func (c *Client) UpdateMealPlanTaskStatus(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) (*types.MealPlanTask, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildChangeMealPlanTaskStatusRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create advanced prep step request")
	}

	var validIngredient *types.MealPlanTask
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating advanced prep step")
	}

	return validIngredient, nil
}

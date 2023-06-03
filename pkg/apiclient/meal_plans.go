package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetMealPlan gets a meal plan.
func (c *Client) GetMealPlan(ctx context.Context, mealPlanID string) (*types.MealPlan, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	req, err := c.requestBuilder.BuildGetMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal plan request")
	}

	var mealPlan *types.MealPlan
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlan); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan")
	}

	return mealPlan, nil
}

// GetMealPlans retrieves a list of meal plans.
func (c *Client) GetMealPlans(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlan], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetMealPlansRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building meal plans list request")
	}

	var mealPlans *types.QueryFilteredResult[types.MealPlan]
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlans); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plans")
	}

	return mealPlans, nil
}

// CreateMealPlan creates a meal plan.
func (c *Client) CreateMealPlan(ctx context.Context, input *types.MealPlanCreationRequestInput) (*types.MealPlan, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create meal plan request")
	}

	var mealPlan *types.MealPlan
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlan); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan")
	}

	return mealPlan, nil
}

// UpdateMealPlan updates a meal plan.
func (c *Client) UpdateMealPlan(ctx context.Context, mealPlan *types.MealPlan) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlan == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlan.ID)
	tracing.AttachMealPlanIDToSpan(span, mealPlan.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanRequest(ctx, mealPlan)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update meal plan request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &mealPlan); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan %s", mealPlan.ID)
	}

	return nil
}

// ArchiveMealPlan archives a meal plan.
func (c *Client) ArchiveMealPlan(ctx context.Context, mealPlanID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	req, err := c.requestBuilder.BuildArchiveMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive meal plan request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan %s", mealPlanID)
	}

	return nil
}

// FinalizeMealPlan gets a meal plan.
func (c *Client) FinalizeMealPlan(ctx context.Context, mealPlanID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	req, err := c.requestBuilder.BuildFinalizeMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building get meal plan request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving meal plan")
	}

	return nil
}

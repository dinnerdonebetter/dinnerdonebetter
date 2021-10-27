package httpclient

import (
	"context"

	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetMealPlan gets a meal plan.
func (c *Client) GetMealPlan(ctx context.Context, mealPlanID string) (*types.MealPlan, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	req, err := c.requestBuilder.BuildGetMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get meal plan request")
	}

	var mealPlan *types.MealPlan
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlan); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving meal plan")
	}

	return mealPlan, nil
}

// GetMealPlans retrieves a list of meal plans.
func (c *Client) GetMealPlans(ctx context.Context, filter *types.QueryFilter) (*types.MealPlanList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetMealPlansRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building meal plans list request")
	}

	var mealPlans *types.MealPlanList
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlans); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving meal plans")
	}

	return mealPlans, nil
}

// CreateMealPlan creates a meal plan.
func (c *Client) CreateMealPlan(ctx context.Context, input *types.MealPlanCreationRequestInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return "", ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return "", observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanRequest(ctx, input)
	if err != nil {
		return "", observability.PrepareError(err, logger, span, "building create meal plan request")
	}

	var pwr *types.PreWriteResponse
	if err = c.fetchAndUnmarshal(ctx, req, &pwr); err != nil {
		return "", observability.PrepareError(err, logger, span, "creating meal plan")
	}

	return pwr.ID, nil
}

// UpdateMealPlan updates a meal plan.
func (c *Client) UpdateMealPlan(ctx context.Context, mealPlan *types.MealPlan) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlan == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlan.ID)
	tracing.AttachMealPlanIDToSpan(span, mealPlan.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanRequest(ctx, mealPlan)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update meal plan request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &mealPlan); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan %s", mealPlan.ID)
	}

	return nil
}

// ArchiveMealPlan archives a meal plan.
func (c *Client) ArchiveMealPlan(ctx context.Context, mealPlanID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	req, err := c.requestBuilder.BuildArchiveMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive meal plan request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving meal plan %s", mealPlanID)
	}

	return nil
}

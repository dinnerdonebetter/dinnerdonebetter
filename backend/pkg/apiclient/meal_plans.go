package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"

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
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	res, err := c.authedGeneratedClient.GetMealPlan(ctx, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get meal plan")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetMealPlans retrieves a list of meal plans.
func (c *Client) GetMealPlans(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlan], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetMealPlansParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetMealPlans(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "meal plans list")
	}

	var apiResponse *types.APIResponse[[]*types.MealPlan]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plans")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.MealPlan]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
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

	body := generated.CreateMealPlanJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateMealPlan(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create meal plan")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlan.ID)

	body := generated.UpdateMealPlanJSONRequestBody{}
	c.copyType(&body, mealPlan)

	res, err := c.authedGeneratedClient.UpdateMealPlan(ctx, mealPlan.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update meal plan")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
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
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	res, err := c.authedGeneratedClient.ArchiveMealPlan(ctx, mealPlanID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive meal plan")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
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
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	res, err := c.authedGeneratedClient.FinalizeMealPlan(ctx, mealPlanID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "get meal plan")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

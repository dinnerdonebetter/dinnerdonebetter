package httpclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetMealPlanOption gets a meal plan option.
func (c *Client) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanEventID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	req, err := c.requestBuilder.BuildGetMealPlanOptionRequest(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building get meal plan option request")
	}

	var mealPlanOption *types.MealPlanOption
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOption); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving meal plan option")
	}

	return mealPlanOption, nil
}

// GetMealPlanOptions retrieves a list of meal plan options.
func (c *Client) GetMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *types.QueryFilter) (*types.MealPlanOptionList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanEventID)

	req, err := c.requestBuilder.BuildGetMealPlanOptionsRequest(ctx, mealPlanID, mealPlanEventID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building meal plan options list request")
	}

	var mealPlanOptions *types.MealPlanOptionList
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOptions); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving meal plan options")
	}

	return mealPlanOptions, nil
}

// CreateMealPlanOption creates a meal plan option.
func (c *Client) CreateMealPlanOption(ctx context.Context, mealPlanEventID string, input *types.MealPlanOptionCreationRequestInput) (*types.MealPlanOption, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanEventID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanOptionRequest(ctx, mealPlanEventID, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building create meal plan option request")
	}

	var mealPlanOption *types.MealPlanOption
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOption); err != nil {
		return nil, observability.PrepareError(err, span, "creating meal plan option")
	}

	return mealPlanOption, nil
}

// UpdateMealPlanOption updates a meal plan option.
func (c *Client) UpdateMealPlanOption(ctx context.Context, mealPlanID string, mealPlanOption *types.MealPlanOption) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOption == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOption.ID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOption.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanOptionRequest(ctx, mealPlanID, mealPlanOption)
	if err != nil {
		return observability.PrepareError(err, span, "building update meal plan option request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOption); err != nil {
		return observability.PrepareError(err, span, "updating meal plan option %s", mealPlanOption.ID)
	}

	return nil
}

// ArchiveMealPlanOption archives a meal plan option.
func (c *Client) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanEventID)

	if mealPlanOptionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	req, err := c.requestBuilder.BuildArchiveMealPlanOptionRequest(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive meal plan option request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "archiving meal plan option %s", mealPlanOptionID)
	}

	return nil
}

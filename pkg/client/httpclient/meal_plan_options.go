package httpclient

import (
	"context"

	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetMealPlanOption gets a meal plan option.
func (c *Client) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	req, err := c.requestBuilder.BuildGetMealPlanOptionRequest(ctx, mealPlanID, mealPlanOptionID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get meal plan option request")
	}

	var mealPlanOption *types.MealPlanOption
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOption); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving meal plan option")
	}

	return mealPlanOption, nil
}

// GetMealPlanOptions retrieves a list of meal plan options.
func (c *Client) GetMealPlanOptions(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (*types.MealPlanOptionList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	req, err := c.requestBuilder.BuildGetMealPlanOptionsRequest(ctx, mealPlanID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building meal plan options list request")
	}

	var mealPlanOptions *types.MealPlanOptionList
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOptions); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving meal plan options")
	}

	return mealPlanOptions, nil
}

// CreateMealPlanOption creates a meal plan option.
func (c *Client) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionCreationRequestInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return "", ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return "", observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanOptionRequest(ctx, input)
	if err != nil {
		return "", observability.PrepareError(err, logger, span, "building create meal plan option request")
	}

	var pwr *types.PreWriteResponse
	if err = c.fetchAndUnmarshal(ctx, req, &pwr); err != nil {
		return "", observability.PrepareError(err, logger, span, "creating meal plan option")
	}

	return pwr.ID, nil
}

// UpdateMealPlanOption updates a meal plan option.
func (c *Client) UpdateMealPlanOption(ctx context.Context, mealPlanOption *types.MealPlanOption) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlanOption == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOption.ID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOption.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanOptionRequest(ctx, mealPlanOption)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update meal plan option request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOption); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option %s", mealPlanOption.ID)
	}

	return nil
}

// ArchiveMealPlanOption archives a meal plan option.
func (c *Client) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	req, err := c.requestBuilder.BuildArchiveMealPlanOptionRequest(ctx, mealPlanID, mealPlanOptionID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive meal plan option request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving meal plan option %s", mealPlanOptionID)
	}

	return nil
}

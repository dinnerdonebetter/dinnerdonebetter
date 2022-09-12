package httpclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetMealPlanEvent gets a meal plan event.
func (c *Client) GetMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
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

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	req, err := c.requestBuilder.BuildGetMealPlanEventRequest(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building get meal plan event request")
	}

	var mealPlanEvent *types.MealPlanEvent
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanEvent); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving meal plan event")
	}

	return mealPlanEvent, nil
}

// GetMealPlanEvents retrieves a list of meal plan events.
func (c *Client) GetMealPlanEvents(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (*types.MealPlanEventList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	req, err := c.requestBuilder.BuildGetMealPlanEventsRequest(ctx, mealPlanID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building meal plan events list request")
	}

	var mealPlanEvents *types.MealPlanEventList
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanEvents); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving meal plan events")
	}

	return mealPlanEvents, nil
}

// CreateMealPlanEvent creates a meal plan event.
func (c *Client) CreateMealPlanEvent(ctx context.Context, mealPlanID string, input *types.MealPlanEventCreationRequestInput) (*types.MealPlanEvent, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanEventRequest(ctx, mealPlanID, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building create meal plan event request")
	}

	var mealPlanEvent *types.MealPlanEvent
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanEvent); err != nil {
		return nil, observability.PrepareError(err, span, "creating meal plan event")
	}

	return mealPlanEvent, nil
}

// UpdateMealPlanEvent updates a meal plan event.
func (c *Client) UpdateMealPlanEvent(ctx context.Context, mealPlanEvent *types.MealPlanEvent) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanEvent == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEvent.ID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEvent.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanEventRequest(ctx, mealPlanEvent)
	if err != nil {
		return observability.PrepareError(err, span, "building update meal plan event request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanEvent); err != nil {
		return observability.PrepareError(err, span, "updating meal plan event %s", mealPlanEvent.ID)
	}

	return nil
}

// ArchiveMealPlanEvent archives a meal plan event.
func (c *Client) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
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

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	req, err := c.requestBuilder.BuildArchiveMealPlanEventRequest(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive meal plan event request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "archiving meal plan event %s", mealPlanEventID)
	}

	return nil
}

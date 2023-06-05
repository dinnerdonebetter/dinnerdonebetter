package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetMealRating gets a meal rating.
func (c *Client) GetMealRating(ctx context.Context, mealID, mealRatingID string) (*types.MealRating, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if mealRatingID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealRatingIDKey, mealRatingID)
	tracing.AttachMealRatingIDToSpan(span, mealRatingID)

	req, err := c.requestBuilder.BuildGetMealRatingRequest(ctx, mealID, mealRatingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal rating request")
	}

	var mealRating *types.MealRating
	if err = c.fetchAndUnmarshal(ctx, req, &mealRating); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal rating")
	}

	return mealRating, nil
}

// GetMealRatings retrieves a list of meal ratings.
func (c *Client) GetMealRatings(ctx context.Context, mealID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealRating], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	req, err := c.requestBuilder.BuildGetMealRatingsRequest(ctx, mealID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building meal ratings list request")
	}

	var mealRatings *types.QueryFilteredResult[types.MealRating]
	if err = c.fetchAndUnmarshal(ctx, req, &mealRatings); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal ratings")
	}

	return mealRatings, nil
}

// CreateMealRating creates a meal rating.
func (c *Client) CreateMealRating(ctx context.Context, mealID string, input *types.MealRatingCreationRequestInput) (*types.MealRating, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealRatingRequest(ctx, mealID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create meal rating request")
	}

	var mealRating *types.MealRating
	if err = c.fetchAndUnmarshal(ctx, req, &mealRating); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal rating")
	}

	return mealRating, nil
}

// UpdateMealRating updates a meal rating.
func (c *Client) UpdateMealRating(ctx context.Context, mealRating *types.MealRating) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealRating == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealRatingIDKey, mealRating.ID)
	tracing.AttachMealRatingIDToSpan(span, mealRating.ID)

	req, err := c.requestBuilder.BuildUpdateMealRatingRequest(ctx, mealRating)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update meal rating request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &mealRating); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal rating %s", mealRating.ID)
	}

	return nil
}

// ArchiveMealRating archives a meal rating.
func (c *Client) ArchiveMealRating(ctx context.Context, mealID, mealRatingID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if mealRatingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealRatingIDKey, mealRatingID)
	tracing.AttachMealRatingIDToSpan(span, mealRatingID)

	req, err := c.requestBuilder.BuildArchiveMealRatingRequest(ctx, mealID, mealRatingID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive meal rating request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal rating %s", mealRatingID)
	}

	return nil
}

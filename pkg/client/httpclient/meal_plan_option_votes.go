package httpclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetMealPlanOptionVote gets a meal plan option vote.
func (c *Client) GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

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

	if mealPlanOptionVoteID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	req, err := c.requestBuilder.BuildGetMealPlanOptionVoteRequest(ctx, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get meal plan option vote request")
	}

	var mealPlanOptionVote *types.MealPlanOptionVote
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOptionVote); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving meal plan option vote")
	}

	return mealPlanOptionVote, nil
}

// GetMealPlanOptionVotes retrieves a list of meal plan option votes.
func (c *Client) GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanOptionID string, filter *types.QueryFilter) (*types.MealPlanOptionVoteList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

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

	req, err := c.requestBuilder.BuildGetMealPlanOptionVotesRequest(ctx, mealPlanID, mealPlanOptionID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building meal plan option votes list request")
	}

	var mealPlanOptionVotes *types.MealPlanOptionVoteList
	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOptionVotes); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving meal plan option votes")
	}

	return mealPlanOptionVotes, nil
}

// CreateMealPlanOptionVote creates a meal plan option vote.
func (c *Client) CreateMealPlanOptionVote(ctx context.Context, mealPlanID string, input *types.MealPlanOptionVoteCreationRequestInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return "", ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if input == nil {
		return "", ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return "", observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanOptionVoteRequest(ctx, mealPlanID, input)
	if err != nil {
		return "", observability.PrepareError(err, logger, span, "building create meal plan option vote request")
	}

	var pwr *types.PreWriteResponse
	if err = c.fetchAndUnmarshal(ctx, req, &pwr); err != nil {
		return "", observability.PrepareError(err, logger, span, "creating meal plan option vote")
	}

	return pwr.ID, nil
}

// UpdateMealPlanOptionVote updates a meal plan option vote.
func (c *Client) UpdateMealPlanOptionVote(ctx context.Context, mealPlanID string, mealPlanOptionVote *types.MealPlanOptionVote) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionVote == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVote.ID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVote.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanOptionVoteRequest(ctx, mealPlanID, mealPlanOptionVote)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update meal plan option vote request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOptionVote); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option vote %s", mealPlanOptionVote.ID)
	}

	return nil
}

// ArchiveMealPlanOptionVote archives a meal plan option vote.
func (c *Client) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

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

	if mealPlanOptionVoteID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	req, err := c.requestBuilder.BuildArchiveMealPlanOptionVoteRequest(ctx, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive meal plan option vote request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving meal plan option vote %s", mealPlanOptionVoteID)
	}

	return nil
}

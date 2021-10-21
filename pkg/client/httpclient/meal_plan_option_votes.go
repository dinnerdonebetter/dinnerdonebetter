package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// GetMealPlanOptionVote gets a meal plan option vote.
func (c *Client) GetMealPlanOptionVote(ctx context.Context, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlanOptionVoteID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	req, err := c.requestBuilder.BuildGetMealPlanOptionVoteRequest(ctx, mealPlanOptionVoteID)
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
func (c *Client) GetMealPlanOptionVotes(ctx context.Context, filter *types.QueryFilter) (*types.MealPlanOptionVoteList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetMealPlanOptionVotesRequest(ctx, filter)
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
func (c *Client) CreateMealPlanOptionVote(ctx context.Context, input *types.MealPlanOptionVoteCreationRequestInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return "", ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return "", observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanOptionVoteRequest(ctx, input)
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
func (c *Client) UpdateMealPlanOptionVote(ctx context.Context, mealPlanOptionVote *types.MealPlanOptionVote) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlanOptionVote == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVote.ID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVote.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanOptionVoteRequest(ctx, mealPlanOptionVote)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update meal plan option vote request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &mealPlanOptionVote); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option vote %s", mealPlanOptionVote.ID)
	}

	return nil
}

// ArchiveMealPlanOptionVote archives a meal plan option vote.
func (c *Client) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanOptionVoteID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if mealPlanOptionVoteID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	req, err := c.requestBuilder.BuildArchiveMealPlanOptionVoteRequest(ctx, mealPlanOptionVoteID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive meal plan option vote request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving meal plan option vote %s", mealPlanOptionVoteID)
	}

	return nil
}

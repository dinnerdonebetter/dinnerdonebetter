package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetMealPlanOptionVote gets a meal plan ClientOption vote.
func (c *Client) GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	if mealPlanOptionVoteID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	req, err := c.requestBuilder.BuildGetMealPlanOptionVoteRequest(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get meal plan ClientOption vote")
	}

	var apiResponse *types.APIResponse[*types.MealPlanOptionVote]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan ClientOption vote")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetMealPlanOptionVotes retrieves a list of meal plan ClientOption votes.
func (c *Client) GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlanOptionVote], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	req, err := c.requestBuilder.BuildGetMealPlanOptionVotesRequest(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "meal plan ClientOption votes list")
	}

	var apiResponse *types.APIResponse[[]*types.MealPlanOptionVote]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan ClientOption votes")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.MealPlanOptionVote]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateMealPlanOptionVote creates a meal plan ClientOption vote.
func (c *Client) CreateMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID string, input *types.MealPlanOptionVoteCreationRequestInput) ([]*types.MealPlanOptionVote, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanOptionVoteRequest(ctx, mealPlanID, mealPlanEventID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create meal plan ClientOption vote")
	}

	var apiResponse *types.APIResponse[[]*types.MealPlanOptionVote]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan ClientOption vote")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateMealPlanOptionVote updates a meal plan ClientOption vote.
func (c *Client) UpdateMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID string, mealPlanOptionVote *types.MealPlanOptionVote) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionVote == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVote.ID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVote.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanOptionVoteRequest(ctx, mealPlanID, mealPlanEventID, mealPlanOptionVote)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update meal plan ClientOption vote")
	}

	var apiResponse *types.APIResponse[*types.MealPlanOptionVote]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan ClientOption vote %s", mealPlanOptionVote.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveMealPlanOptionVote archives a meal plan ClientOption vote.
func (c *Client) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	if mealPlanOptionVoteID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	req, err := c.requestBuilder.BuildArchiveMealPlanOptionVoteRequest(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive meal plan ClientOption vote")
	}

	var apiResponse *types.APIResponse[*types.MealPlanOptionVote]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan ClientOption vote %s", mealPlanOptionVoteID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

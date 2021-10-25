package requests

import (
	"context"
	"net/http"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	mealPlanOptionVotesBasePath = "meal_plan_option_votes"
)

// BuildGetMealPlanOptionVoteRequest builds an HTTP request for fetching a meal plan option vote.
func (b *Builder) BuildGetMealPlanOptionVoteRequest(ctx context.Context, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

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

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
		mealPlanOptionVotesBasePath,
		mealPlanOptionVoteID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetMealPlanOptionVotesRequest builds an HTTP request for fetching a list of meal plan option votes.
func (b *Builder) BuildGetMealPlanOptionVotesRequest(ctx context.Context, mealPlanID, mealPlanOptionID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

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

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		mealPlansBasePath,
		mealPlanID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
		mealPlanOptionVotesBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateMealPlanOptionVoteRequest builds an HTTP request for creating a meal plan option vote.
func (b *Builder) BuildCreateMealPlanOptionVoteRequest(ctx context.Context, mealPlanID string, input *types.MealPlanOptionVoteCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanOptionsBasePath,
		input.BelongsToMealPlanOption,
		mealPlanOptionVotesBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateMealPlanOptionVoteRequest builds an HTTP request for updating a meal plan option vote.
func (b *Builder) BuildUpdateMealPlanOptionVoteRequest(ctx context.Context, mealPlanID string, mealPlanOptionVote *types.MealPlanOptionVote) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionVote == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVote.ID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVote.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanOptionsBasePath,
		mealPlanOptionVote.BelongsToMealPlanOption,
		mealPlanOptionVotesBasePath,
		mealPlanOptionVote.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, mealPlanOptionVote)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveMealPlanOptionVoteRequest builds an HTTP request for archiving a meal plan option vote.
func (b *Builder) BuildArchiveMealPlanOptionVoteRequest(ctx context.Context, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

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

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
		mealPlanOptionVotesBasePath,
		mealPlanOptionVoteID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	mealPlanOptionVotesBasePath = "meal_plan_option_votes"
)

// BuildGetMealPlanOptionVoteRequest builds an HTTP request for fetching a meal plan option vote.
func (b *Builder) BuildGetMealPlanOptionVoteRequest(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	if mealPlanOptionVoteID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanEventID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
		mealPlanOptionVotesBasePath,
		mealPlanOptionVoteID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetMealPlanOptionVotesRequest builds an HTTP request for fetching a list of meal plan option votes.
func (b *Builder) BuildGetMealPlanOptionVotesRequest(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanEventID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
		mealPlanOptionVotesBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateMealPlanOptionVoteRequest builds an HTTP request for creating a meal plan option vote.
func (b *Builder) BuildCreateMealPlanOptionVoteRequest(ctx context.Context, mealPlanID, mealPlanEventID string, input *types.MealPlanOptionVoteCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanEventID,
		"vote",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateMealPlanOptionVoteRequest builds an HTTP request for updating a meal plan option vote.
func (b *Builder) BuildUpdateMealPlanOptionVoteRequest(ctx context.Context, mealPlanID, mealPlanEventID string, mealPlanOptionVote *types.MealPlanOptionVote) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	if mealPlanOptionVote == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVote.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanEventID,
		mealPlanOptionsBasePath,
		mealPlanOptionVote.BelongsToMealPlanOption,
		mealPlanOptionVotesBasePath,
		mealPlanOptionVote.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := types.MealPlanOptionVoteUpdateRequestInputFromMealPlanOptionVote(mealPlanOptionVote)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveMealPlanOptionVoteRequest builds an HTTP request for archiving a meal plan option vote.
func (b *Builder) BuildArchiveMealPlanOptionVoteRequest(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	if mealPlanOptionVoteID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanEventID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
		mealPlanOptionVotesBasePath,
		mealPlanOptionVoteID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

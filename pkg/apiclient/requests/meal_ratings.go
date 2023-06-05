package requests

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	mealRatingsBasePath = "ratings"
)

// BuildGetMealRatingRequest builds an HTTP request for fetching a valid instrument.
func (b *Builder) BuildGetMealRatingRequest(ctx context.Context, mealID, mealRatingID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealIDToSpan(span, mealID)

	if mealRatingID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealRatingIDToSpan(span, mealRatingID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealsBasePath,
		mealID,
		mealRatingsBasePath,
		mealRatingID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetMealRatingsRequest builds an HTTP request for fetching a list of valid instruments.
func (b *Builder) BuildGetMealRatingsRequest(ctx context.Context, mealID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealIDToSpan(span, mealID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		mealsBasePath,
		mealID,
		mealRatingsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateMealRatingRequest builds an HTTP request for creating a valid instrument.
func (b *Builder) BuildCreateMealRatingRequest(ctx context.Context, mealID string, input *types.MealRatingCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealIDToSpan(span, mealID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(
		ctx,
		nil,
		mealsBasePath,
		mealID,
		mealRatingsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateMealRatingRequest builds an HTTP request for updating a valid instrument.
func (b *Builder) BuildUpdateMealRatingRequest(ctx context.Context, mealRating *types.MealRating) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealRating == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachMealRatingIDToSpan(span, mealRating.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealsBasePath,
		mealRating.MealID,
		mealRatingsBasePath,
		mealRating.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertMealRatingToMealRatingUpdateRequestInput(mealRating)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveMealRatingRequest builds an HTTP request for archiving a valid instrument.
func (b *Builder) BuildArchiveMealRatingRequest(ctx context.Context, mealID, mealRatingID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealIDToSpan(span, mealID)

	if mealRatingID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealRatingIDToSpan(span, mealRatingID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealsBasePath,
		mealID,
		mealRatingsBasePath,
		mealRatingID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

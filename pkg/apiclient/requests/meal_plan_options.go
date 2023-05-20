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
	mealPlanOptionsBasePath = "options"
)

// BuildGetMealPlanOptionRequest builds an HTTP request for fetching a meal plan option.
func (b *Builder) BuildGetMealPlanOptionRequest(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*http.Request, error) {
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
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanEventID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetMealPlanOptionsRequest builds an HTTP request for fetching a list of meal plan options.
func (b *Builder) BuildGetMealPlanOptionsRequest(ctx context.Context, mealPlanID, mealPlanEventID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanEventID,
		mealPlanOptionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateMealPlanOptionRequest builds an HTTP request for creating a meal plan option.
func (b *Builder) BuildCreateMealPlanOptionRequest(ctx context.Context, mealPlanID, mealPlanEventID string, input *types.MealPlanOptionCreationRequestInput) (*http.Request, error) {
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
		mealPlanOptionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateMealPlanOptionRequest builds an HTTP request for updating a meal plan option.
func (b *Builder) BuildUpdateMealPlanOptionRequest(ctx context.Context, mealPlanID string, mealPlanOption *types.MealPlanOption) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanOption == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOption.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanOption.BelongsToMealPlanEvent,
		mealPlanOptionsBasePath,
		mealPlanOption.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput(mealPlanOption)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveMealPlanOptionRequest builds an HTTP request for archiving a meal plan option.
func (b *Builder) BuildArchiveMealPlanOptionRequest(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanEventsBasePath,
		mealPlanEventID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

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
	mealPlansBasePath = "meal_plans"
)

// BuildGetMealPlanRequest builds an HTTP request for fetching a meal plan.
func (b *Builder) BuildGetMealPlanRequest(ctx context.Context, mealPlanID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetMealPlansRequest builds an HTTP request for fetching a list of meal plans.
func (b *Builder) BuildGetMealPlansRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		mealPlansBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateMealPlanRequest builds an HTTP request for creating a meal plan.
func (b *Builder) BuildCreateMealPlanRequest(ctx context.Context, input *types.MealPlanCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

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
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateMealPlanRequest builds an HTTP request for updating a meal plan.
func (b *Builder) BuildUpdateMealPlanRequest(ctx context.Context, mealPlan *types.MealPlan) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlan == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlan.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlan.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertMealPlanToMealPlanUpdateRequestInput(mealPlan)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveMealPlanRequest builds an HTTP request for archiving a meal plan.
func (b *Builder) BuildArchiveMealPlanRequest(ctx context.Context, mealPlanID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

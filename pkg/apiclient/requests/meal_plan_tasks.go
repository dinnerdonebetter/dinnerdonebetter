package requests

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	mealPlanTasksBasePath = "tasks"
)

// BuildGetMealPlanTaskRequest builds an HTTP request for fetching a meal plan.
func (b *Builder) BuildGetMealPlanTaskRequest(ctx context.Context, mealPlanID, mealPlanTaskID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanTaskIDToSpan(span, mealPlanTaskID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanTasksBasePath,
		mealPlanTaskID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateMealPlanTaskRequest builds an HTTP request for fetching a meal plan.
func (b *Builder) BuildCreateMealPlanTaskRequest(ctx context.Context, mealPlanID string, input *types.MealPlanTaskCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanTasksBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetMealPlanTasksRequest builds an HTTP request for fetching a list of meal plan tasks.
func (b *Builder) BuildGetMealPlanTasksRequest(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (*http.Request, error) {
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
		mealPlanTasksBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildChangeMealPlanTaskStatusRequest builds an HTTP request for archiving a meal plan.
func (b *Builder) BuildChangeMealPlanTaskStatusRequest(ctx context.Context, mealPlanID string, input *types.MealPlanTaskStatusChangeRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachMealPlanTaskIDToSpan(span, input.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanTasksBasePath,
		input.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPatch, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

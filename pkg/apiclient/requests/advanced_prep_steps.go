package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	mealPlanTasksBasePath = "meal_plan_tasks"
)

// BuildGetMealPlanTaskRequest builds an HTTP request for fetching a meal plan.
func (b *Builder) BuildGetMealPlanTaskRequest(ctx context.Context, mealPlanTaskID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanTaskIDToSpan(span, mealPlanTaskID)

	uri := b.BuildURL(
		ctx,
		nil,
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

// BuildGetMealPlanTasksRequest builds an HTTP request for fetching a list of advanced prep steps.
func (b *Builder) BuildGetMealPlanTasksRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
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
func (b *Builder) BuildChangeMealPlanTaskStatusRequest(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachMealPlanTaskIDToSpan(span, input.ID)

	uri := b.BuildURL(
		ctx,
		nil,
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

package requests

import (
	"context"
	"net/http"

	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	mealPlanOptionsBasePath = "meal_plan_options"
)

// BuildGetMealPlanOptionRequest builds an HTTP request for fetching a meal plan option.
func (b *Builder) BuildGetMealPlanOptionRequest(ctx context.Context, mealPlanID, mealPlanOptionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

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
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetMealPlanOptionsRequest builds an HTTP request for fetching a list of meal plan options.
func (b *Builder) BuildGetMealPlanOptionsRequest(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		mealPlansBasePath,
		mealPlanID,
		mealPlanOptionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateMealPlanOptionRequest builds an HTTP request for creating a meal plan option.
func (b *Builder) BuildCreateMealPlanOptionRequest(ctx context.Context, input *types.MealPlanOptionCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

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
		input.BelongsToMealPlan,
		mealPlanOptionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateMealPlanOptionRequest builds an HTTP request for updating a meal plan option.
func (b *Builder) BuildUpdateMealPlanOptionRequest(ctx context.Context, mealPlanOption *types.MealPlanOption) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if mealPlanOption == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOption.ID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOption.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanOption.BelongsToMealPlan,
		mealPlanOptionsBasePath,
		mealPlanOption.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, mealPlanOption)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveMealPlanOptionRequest builds an HTTP request for archiving a meal plan option.
func (b *Builder) BuildArchiveMealPlanOptionRequest(ctx context.Context, mealPlanID, mealPlanOptionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

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
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanOptionsBasePath,
		mealPlanOptionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

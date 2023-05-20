package requests

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	mealPlanGroceryListItemsBasePath = "grocery_list_items"
)

// BuildGetMealPlanGroceryListItemRequest builds an HTTP request for fetching a meal plan grocery list item.
func (b *Builder) BuildGetMealPlanGroceryListItemRequest(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItemID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanGroceryListItemsBasePath,
		mealPlanGroceryListItemID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateMealPlanGroceryListItemRequest builds an HTTP request for fetching a meal plan grocery list item.
func (b *Builder) BuildCreateMealPlanGroceryListItemRequest(ctx context.Context, mealPlanID string, input *types.MealPlanGroceryListItemCreationRequestInput) (*http.Request, error) {
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
		mealPlanGroceryListItemsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetMealPlanGroceryListItemsForMealPlanRequest builds an HTTP request for fetching a list of meal plan grocery list items.
func (b *Builder) BuildGetMealPlanGroceryListItemsForMealPlanRequest(ctx context.Context, mealPlanID string) (*http.Request, error) {
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
		mealPlanGroceryListItemsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateMealPlanGroceryListItemRequest builds an HTTP request for updating a meal plan grocery list item.
func (b *Builder) BuildUpdateMealPlanGroceryListItemRequest(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string, input *types.MealPlanGroceryListItemUpdateRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItemID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanGroceryListItemsBasePath,
		mealPlanGroceryListItemID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPatch, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveMealPlanGroceryListItemRequest builds an HTTP request for archiving a meal plan grocery list item.
func (b *Builder) BuildArchiveMealPlanGroceryListItemRequest(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, mealPlanGroceryListItemID)

	uri := b.BuildURL(
		ctx,
		nil,
		mealPlansBasePath,
		mealPlanID,
		mealPlanGroceryListItemsBasePath,
		mealPlanGroceryListItemID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

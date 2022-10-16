package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
)

const (
	recipePrepTasksBasePath = "prep_tasks"
)

// BuildGetRecipePrepTaskRequest builds an HTTP request for fetching a recipe step.
func (b *Builder) BuildGetRecipePrepTaskRequest(ctx context.Context, recipeID, recipePrepTaskID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipePrepTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipePrepTasksBasePath,
		recipePrepTaskID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRecipePrepTasksRequest builds an HTTP request for fetching a list of recipe steps.
func (b *Builder) BuildGetRecipePrepTasksRequest(ctx context.Context, recipeID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		recipesBasePath,
		recipeID,
		recipePrepTasksBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateRecipePrepTaskRequest builds an HTTP request for creating a recipe step.
func (b *Builder) BuildCreateRecipePrepTaskRequest(ctx context.Context, input *types.RecipePrepTaskCreationRequestInput) (*http.Request, error) {
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
		recipesBasePath,
		input.BelongsToRecipe,
		recipePrepTasksBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateRecipePrepTaskRequest builds an HTTP request for updating a recipe step.
func (b *Builder) BuildUpdateRecipePrepTaskRequest(ctx context.Context, recipePrepTask *types.RecipePrepTask) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipePrepTask == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTask.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipePrepTask.BelongsToRecipe,
		recipePrepTasksBasePath,
		recipePrepTask.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput(recipePrepTask)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveRecipePrepTaskRequest builds an HTTP request for archiving a recipe step.
func (b *Builder) BuildArchiveRecipePrepTaskRequest(ctx context.Context, recipeID, recipePrepTaskID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipePrepTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipePrepTasksBasePath,
		recipePrepTaskID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

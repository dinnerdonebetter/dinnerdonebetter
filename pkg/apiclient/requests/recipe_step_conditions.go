package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

const (
	recipeStepConditionsBasePath = "conditions"
)

// BuildGetRecipeStepConditionRequest builds an HTTP request for fetching a recipe step condition.
func (b *Builder) BuildGetRecipeStepConditionRequest(ctx context.Context, recipeID, recipeStepID, recipeStepConditionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepConditionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepConditionID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepConditionsBasePath,
		recipeStepConditionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRecipeStepConditionsRequest builds an HTTP request for fetching a list of recipe step conditions.
func (b *Builder) BuildGetRecipeStepConditionsRequest(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepConditionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateRecipeStepConditionRequest builds an HTTP request for creating a recipe step condition.
func (b *Builder) BuildCreateRecipeStepConditionRequest(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepConditionCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrEmptyInputProvided
	}
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

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
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepConditionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateRecipeStepConditionRequest builds an HTTP request for updating a recipe step condition.
func (b *Builder) BuildUpdateRecipeStepConditionRequest(ctx context.Context, recipeID string, recipeStepCondition *types.RecipeStepCondition) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepCondition == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepCondition.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepCondition.BelongsToRecipeStep,
		recipeStepConditionsBasePath,
		recipeStepCondition.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertRecipeStepConditionToRecipeStepConditionUpdateRequestInput(recipeStepCondition)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveRecipeStepConditionRequest builds an HTTP request for archiving a recipe step condition.
func (b *Builder) BuildArchiveRecipeStepConditionRequest(ctx context.Context, recipeID, recipeStepID, recipeStepConditionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepConditionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepConditionID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepConditionsBasePath,
		recipeStepConditionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

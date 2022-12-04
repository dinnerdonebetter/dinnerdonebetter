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
	recipeStepCompletionConditionsBasePath = "completion_conditions"
)

// BuildGetRecipeStepCompletionConditionRequest builds an HTTP request for fetching a recipe step completion condition.
func (b *Builder) BuildGetRecipeStepCompletionConditionRequest(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*http.Request, error) {
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

	if recipeStepCompletionConditionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionConditionID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepCompletionConditionsBasePath,
		recipeStepCompletionConditionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRecipeStepCompletionConditionsRequest builds an HTTP request for fetching a list of recipe step completion conditions.
func (b *Builder) BuildGetRecipeStepCompletionConditionsRequest(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*http.Request, error) {
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
		recipeStepCompletionConditionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateRecipeStepCompletionConditionRequest builds an HTTP request for creating a recipe step completion condition.
func (b *Builder) BuildCreateRecipeStepCompletionConditionRequest(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionCreationRequestInput) (*http.Request, error) {
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
		recipeStepCompletionConditionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateRecipeStepCompletionConditionRequest builds an HTTP request for updating a recipe step completion condition.
func (b *Builder) BuildUpdateRecipeStepCompletionConditionRequest(ctx context.Context, recipeID string, recipeStepCompletionCondition *types.RecipeStepCompletionCondition) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepCompletionCondition == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionCondition.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepCompletionCondition.BelongsToRecipeStep,
		recipeStepCompletionConditionsBasePath,
		recipeStepCompletionCondition.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionUpdateRequestInput(recipeStepCompletionCondition)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveRecipeStepCompletionConditionRequest builds an HTTP request for archiving a recipe step completion condition.
func (b *Builder) BuildArchiveRecipeStepCompletionConditionRequest(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*http.Request, error) {
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

	if recipeStepCompletionConditionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionConditionID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepCompletionConditionsBasePath,
		recipeStepCompletionConditionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

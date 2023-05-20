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
	recipeStepInstrumentsBasePath = "instruments"
)

// BuildGetRecipeStepInstrumentRequest builds an HTTP request for fetching a recipe step instrument.
func (b *Builder) BuildGetRecipeStepInstrumentRequest(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*http.Request, error) {
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

	if recipeStepInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepInstrumentsBasePath,
		recipeStepInstrumentID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRecipeStepInstrumentsRequest builds an HTTP request for fetching a list of recipe step instruments.
func (b *Builder) BuildGetRecipeStepInstrumentsRequest(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*http.Request, error) {
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
		recipeStepInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateRecipeStepInstrumentRequest builds an HTTP request for creating a recipe step instrument.
func (b *Builder) BuildCreateRecipeStepInstrumentRequest(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentCreationRequestInput) (*http.Request, error) {
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
		recipeStepInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateRecipeStepInstrumentRequest builds an HTTP request for updating a recipe step instrument.
func (b *Builder) BuildUpdateRecipeStepInstrumentRequest(ctx context.Context, recipeID string, recipeStepInstrument *types.RecipeStepInstrument) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepInstrument == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrument.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepInstrument.BelongsToRecipeStep,
		recipeStepInstrumentsBasePath,
		recipeStepInstrument.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(recipeStepInstrument)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveRecipeStepInstrumentRequest builds an HTTP request for archiving a recipe step instrument.
func (b *Builder) BuildArchiveRecipeStepInstrumentRequest(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*http.Request, error) {
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

	if recipeStepInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		recipeStepsBasePath,
		recipeStepID,
		recipeStepInstrumentsBasePath,
		recipeStepInstrumentID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

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
	validIngredientPreparationsBasePath = "valid_ingredient_preparations"
)

// BuildGetValidIngredientPreparationRequest builds an HTTP request for fetching a valid ingredient preparation.
func (b *Builder) BuildGetValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparationID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientPreparationsBasePath,
		validIngredientPreparationID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientPreparationsRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidIngredientPreparationsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientPreparationsForIngredientRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidIngredientPreparationsForIngredientRequest(ctx context.Context, ingredientID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, ingredientID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientPreparationsBasePath,
		"by_ingredient",
		ingredientID,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientPreparationsForPreparationRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidIngredientPreparationsForPreparationRequest(ctx context.Context, preparationID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, preparationID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientPreparationsBasePath,
		"by_preparation",
		preparationID,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidIngredientPreparationRequest builds an HTTP request for creating a valid ingredient preparation.
func (b *Builder) BuildCreateValidIngredientPreparationRequest(ctx context.Context, input *types.ValidIngredientPreparationCreationRequestInput) (*http.Request, error) {
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
		validIngredientPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidIngredientPreparationRequest builds an HTTP request for updating a valid ingredient preparation.
func (b *Builder) BuildUpdateValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparation *types.ValidIngredientPreparation) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientPreparation == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparation.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientPreparationsBasePath,
		validIngredientPreparation.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertValidIngredientPreparationToValidIngredientPreparationUpdateRequestInput(validIngredientPreparation)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidIngredientPreparationRequest builds an HTTP request for archiving a valid ingredient preparation.
func (b *Builder) BuildArchiveValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparationID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientPreparationsBasePath,
		validIngredientPreparationID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

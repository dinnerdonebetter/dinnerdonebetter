package requests

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

const (
	validIngredientsBasePath = "valid_ingredients"
)

// BuildGetValidIngredientRequest builds an HTTP request for fetching a valid ingredient.
func (b *Builder) BuildGetValidIngredientRequest(ctx context.Context, validIngredientID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientsBasePath,
		validIngredientID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRandomValidIngredientRequest builds an HTTP request for fetching a valid ingredient.
func (b *Builder) BuildGetRandomValidIngredientRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientsBasePath,
		randomBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildSearchValidIngredientsRequest builds an HTTP request for querying valid ingredients.
func (b *Builder) BuildSearchValidIngredientsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	params := url.Values{}
	params.Set(types.SearchQueryKey, query)
	params.Set(types.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := b.BuildURL(
		ctx,
		params,
		validIngredientsBasePath,
		"search",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientsRequest builds an HTTP request for fetching a list of valid ingredients.
func (b *Builder) BuildGetValidIngredientsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidIngredientRequest builds an HTTP request for creating a valid ingredient.
func (b *Builder) BuildCreateValidIngredientRequest(ctx context.Context, input *types.ValidIngredientCreationRequestInput) (*http.Request, error) {
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
		validIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidIngredientRequest builds an HTTP request for updating a valid ingredient.
func (b *Builder) BuildUpdateValidIngredientRequest(ctx context.Context, validIngredient *types.ValidIngredient) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredient == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, validIngredient.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientsBasePath,
		validIngredient.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertValidIngredientToValidIngredientUpdateRequestInput(validIngredient)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidIngredientRequest builds an HTTP request for archiving a valid ingredient.
func (b *Builder) BuildArchiveValidIngredientRequest(ctx context.Context, validIngredientID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientsBasePath,
		validIngredientID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

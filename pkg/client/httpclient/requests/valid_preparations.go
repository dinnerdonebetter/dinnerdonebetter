package requests

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	validPreparationsBasePath = "valid_preparations"
)

// BuildGetValidPreparationRequest builds an HTTP request for fetching a valid preparation.
func (b *Builder) BuildGetValidPreparationRequest(ctx context.Context, validPreparationID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		validPreparationID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetRandomValidPreparationRequest builds an HTTP request for fetching a valid preparation.
func (b *Builder) BuildGetRandomValidPreparationRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		randomBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildSearchValidPreparationsRequest builds an HTTP request for querying valid preparations.
func (b *Builder) BuildSearchValidPreparationsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.WithValue(types.SearchQueryKey, query).WithValue(types.LimitQueryKey, limit)

	params := url.Values{}
	params.Set(types.SearchQueryKey, query)
	params.Set(types.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := b.BuildURL(
		ctx,
		params,
		validPreparationsBasePath,
		"search",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidPreparationsRequest builds an HTTP request for fetching a list of valid preparations.
func (b *Builder) BuildGetValidPreparationsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateValidPreparationRequest builds an HTTP request for creating a valid preparation.
func (b *Builder) BuildCreateValidPreparationRequest(ctx context.Context, input *types.ValidPreparationCreationRequestInput) (*http.Request, error) {
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
		validPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidPreparationRequest builds an HTTP request for updating a valid preparation.
func (b *Builder) BuildUpdateValidPreparationRequest(ctx context.Context, validPreparation *types.ValidPreparation) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validPreparation == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparation.ID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparation.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		validPreparation.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, validPreparation)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidPreparationRequest builds an HTTP request for archiving a valid preparation.
func (b *Builder) BuildArchiveValidPreparationRequest(ctx context.Context, validPreparationID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		validPreparationID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

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
	validInstrumentsBasePath = "valid_instruments"
)

// BuildGetValidInstrumentRequest builds an HTTP request for fetching a valid instrument.
func (b *Builder) BuildGetValidInstrumentRequest(ctx context.Context, validInstrumentID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validInstrumentsBasePath,
		validInstrumentID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetRandomValidInstrumentRequest builds an HTTP request for fetching a valid instrument.
func (b *Builder) BuildGetRandomValidInstrumentRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	uri := b.BuildURL(
		ctx,
		nil,
		validInstrumentsBasePath,
		randomBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildSearchValidInstrumentsRequest builds an HTTP request for querying valid instruments.
func (b *Builder) BuildSearchValidInstrumentsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.WithValue(types.SearchQueryKey, query).WithValue(types.LimitQueryKey, limit)

	params := url.Values{}
	params.Set(types.SearchQueryKey, query)
	params.Set(types.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := b.BuildURL(
		ctx,
		params,
		validInstrumentsBasePath,
		"search",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidInstrumentsRequest builds an HTTP request for fetching a list of valid instruments.
func (b *Builder) BuildGetValidInstrumentsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateValidInstrumentRequest builds an HTTP request for creating a valid instrument.
func (b *Builder) BuildCreateValidInstrumentRequest(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (*http.Request, error) {
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
		validInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidInstrumentRequest builds an HTTP request for updating a valid instrument.
func (b *Builder) BuildUpdateValidInstrumentRequest(ctx context.Context, validInstrument *types.ValidInstrument) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validInstrument == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrument.ID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrument.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validInstrumentsBasePath,
		validInstrument.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, validInstrument)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidInstrumentRequest builds an HTTP request for archiving a valid instrument.
func (b *Builder) BuildArchiveValidInstrumentRequest(ctx context.Context, validInstrumentID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validInstrumentsBasePath,
		validInstrumentID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

package requests

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	validPreparationsBasePath = "valid_preparations"
)

// BuildValidPreparationExistsRequest builds an HTTP request for checking the existence of a valid preparation.
func (b *Builder) BuildValidPreparationExistsRequest(ctx context.Context, validPreparationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		id(validPreparationID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidPreparationRequest builds an HTTP request for fetching a valid preparation.
func (b *Builder) BuildGetValidPreparationRequest(ctx context.Context, validPreparationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		id(validPreparationID),
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

	logger := b.logger

	logger = logger.WithValue(types.SearchQueryKey, query).WithValue(types.LimitQueryKey, limit)

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
func (b *Builder) BuildCreateValidPreparationRequest(ctx context.Context, input *types.ValidPreparationCreationInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

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

	logger := b.logger

	if validPreparation == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparation.ID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparation.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(validPreparation.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, validPreparation)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidPreparationRequest builds an HTTP request for archiving a valid preparation.
func (b *Builder) BuildArchiveValidPreparationRequest(ctx context.Context, validPreparationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		id(validPreparationID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetAuditLogForValidPreparationRequest builds an HTTP request for fetching a list of audit log entries pertaining to a valid preparation.
func (b *Builder) BuildGetAuditLogForValidPreparationRequest(ctx context.Context, validPreparationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationsBasePath,
		id(validPreparationID),
		"audit",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

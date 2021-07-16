package requests

import (
	"context"
	"net/http"
	"strconv"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	validPreparationInstrumentsBasePath = "valid_preparation_instruments"
)

// BuildValidPreparationInstrumentExistsRequest builds an HTTP request for checking the existence of a valid preparation instrument.
func (b *Builder) BuildValidPreparationInstrumentExistsRequest(ctx context.Context, validPreparationInstrumentID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		id(validPreparationInstrumentID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidPreparationInstrumentRequest builds an HTTP request for fetching a valid preparation instrument.
func (b *Builder) BuildGetValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrumentID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		id(validPreparationInstrumentID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidPreparationInstrumentsRequest builds an HTTP request for fetching a list of valid preparation instruments.
func (b *Builder) BuildGetValidPreparationInstrumentsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateValidPreparationInstrumentRequest builds an HTTP request for creating a valid preparation instrument.
func (b *Builder) BuildCreateValidPreparationInstrumentRequest(ctx context.Context, input *types.ValidPreparationInstrumentCreationInput) (*http.Request, error) {
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
		validPreparationInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidPreparationInstrumentRequest builds an HTTP request for updating a valid preparation instrument.
func (b *Builder) BuildUpdateValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrument *types.ValidPreparationInstrument) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationInstrument == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrument.ID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrument.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		strconv.FormatUint(validPreparationInstrument.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, validPreparationInstrument)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidPreparationInstrumentRequest builds an HTTP request for archiving a valid preparation instrument.
func (b *Builder) BuildArchiveValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrumentID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		id(validPreparationInstrumentID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetAuditLogForValidPreparationInstrumentRequest builds an HTTP request for fetching a list of audit log entries pertaining to a valid preparation instrument.
func (b *Builder) BuildGetAuditLogForValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrumentID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validPreparationInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		id(validPreparationInstrumentID),
		"audit",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

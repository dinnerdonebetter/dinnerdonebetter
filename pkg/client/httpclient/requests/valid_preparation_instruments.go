package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	validPreparationInstrumentsBasePath = "valid_preparation_instruments"
)

// BuildGetValidPreparationInstrumentRequest builds an HTTP request for fetching a valid ingredient preparation.
func (b *Builder) BuildGetValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrumentID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		validPreparationInstrumentID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidPreparationInstrumentsRequest builds an HTTP request for fetching a list of valid ingredient preparations.
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidPreparationInstrumentsForPreparationRequest builds an HTTP request for fetching a list of valid preparation instruments.
func (b *Builder) BuildGetValidPreparationInstrumentsForPreparationRequest(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidIngredientIDToSpan(span, validPreparationID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationInstrumentsBasePath,
		"by_preparation",
		validPreparationID,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidPreparationInstrumentsForInstrumentRequest builds an HTTP request for fetching a list of valid preparation instruments.
func (b *Builder) BuildGetValidPreparationInstrumentsForInstrumentRequest(ctx context.Context, validInstrumentID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validInstrumentID)
	tracing.AttachValidIngredientIDToSpan(span, validInstrumentID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationInstrumentsBasePath,
		"by_instrument",
		validInstrumentID,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateValidPreparationInstrumentRequest builds an HTTP request for creating a valid ingredient preparation.
func (b *Builder) BuildCreateValidPreparationInstrumentRequest(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*http.Request, error) {
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
		validPreparationInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidPreparationInstrumentRequest builds an HTTP request for updating a valid ingredient preparation.
func (b *Builder) BuildUpdateValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrument *types.ValidPreparationInstrument) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validPreparationInstrument == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrument.ID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrument.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		validPreparationInstrument.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := types.ValidPreparationInstrumentFromValidPreparationInstrument(validPreparationInstrument)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidPreparationInstrumentRequest builds an HTTP request for archiving a valid ingredient preparation.
func (b *Builder) BuildArchiveValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrumentID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.Clone()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		validPreparationInstrumentID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

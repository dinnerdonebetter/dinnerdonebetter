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
	validIngredientPreparationsBasePath = "valid_ingredient_preparations"
)

// BuildValidIngredientPreparationExistsRequest builds an HTTP request for checking the existence of a valid ingredient preparation.
func (b *Builder) BuildValidIngredientPreparationExistsRequest(ctx context.Context, validIngredientPreparationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validIngredientPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientPreparationsBasePath,
		id(validIngredientPreparationID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidIngredientPreparationRequest builds an HTTP request for fetching a valid ingredient preparation.
func (b *Builder) BuildGetValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validIngredientPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientPreparationsBasePath,
		id(validIngredientPreparationID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetValidIngredientPreparationsRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidIngredientPreparationsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateValidIngredientPreparationRequest builds an HTTP request for creating a valid ingredient preparation.
func (b *Builder) BuildCreateValidIngredientPreparationRequest(ctx context.Context, input *types.ValidIngredientPreparationCreationInput) (*http.Request, error) {
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
		validIngredientPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidIngredientPreparationRequest builds an HTTP request for updating a valid ingredient preparation.
func (b *Builder) BuildUpdateValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparation *types.ValidIngredientPreparation) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validIngredientPreparation == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparation.ID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparation.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientPreparationsBasePath,
		strconv.FormatUint(validIngredientPreparation.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, validIngredientPreparation)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidIngredientPreparationRequest builds an HTTP request for archiving a valid ingredient preparation.
func (b *Builder) BuildArchiveValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validIngredientPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientPreparationsBasePath,
		id(validIngredientPreparationID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetAuditLogForValidIngredientPreparationRequest builds an HTTP request for fetching a list of audit log entries pertaining to a valid ingredient preparation.
func (b *Builder) BuildGetAuditLogForValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if validIngredientPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientPreparationsBasePath,
		id(validIngredientPreparationID),
		"audit",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

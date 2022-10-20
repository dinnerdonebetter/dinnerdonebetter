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
	validMeasurementConversionsBasePath = "valid_measurement_conversions"
)

// BuildGetValidMeasurementConversionRequest builds an HTTP request for fetching a valid measurement conversion.
func (b *Builder) BuildGetValidMeasurementConversionRequest(ctx context.Context, validMeasurementConversionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)

	uri := b.BuildURL(
		ctx,
		nil,
		validMeasurementConversionsBasePath,
		validMeasurementConversionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidMeasurementConversionRequest builds an HTTP request for creating a valid measurement conversion.
func (b *Builder) BuildCreateValidMeasurementConversionRequest(ctx context.Context, input *types.ValidMeasurementConversionCreationRequestInput) (*http.Request, error) {
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
		validMeasurementConversionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidMeasurementConversionRequest builds an HTTP request for updating a valid measurement conversion.
func (b *Builder) BuildUpdateValidMeasurementConversionRequest(ctx context.Context, validMeasurementConversion *types.ValidMeasurementConversion) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementConversion == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversion.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validMeasurementConversionsBasePath,
		validMeasurementConversion.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertValidMeasurementConversionToValidMeasurementConversionUpdateRequestInput(validMeasurementConversion)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidMeasurementConversionRequest builds an HTTP request for archiving a valid measurement conversion.
func (b *Builder) BuildArchiveValidMeasurementConversionRequest(ctx context.Context, validMeasurementConversionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)

	uri := b.BuildURL(
		ctx,
		nil,
		validMeasurementConversionsBasePath,
		validMeasurementConversionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

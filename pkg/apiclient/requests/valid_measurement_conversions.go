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
	validMeasurementUnitConversionsBasePath = "valid_measurement_conversions"
)

// BuildGetValidMeasurementUnitConversionRequest builds an HTTP request for fetching a valid measurement conversion.
func (b *Builder) BuildGetValidMeasurementUnitConversionRequest(ctx context.Context, validMeasurementUnitConversionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementUnitConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	uri := b.BuildURL(
		ctx,
		nil,
		validMeasurementUnitConversionsBasePath,
		validMeasurementUnitConversionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidMeasurementUnitConversionsFromUnitRequest builds an HTTP request for fetching a valid measurement conversion.
func (b *Builder) BuildGetValidMeasurementUnitConversionsFromUnitRequest(ctx context.Context, validMeasurementUnitID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	uri := b.BuildURL(
		ctx,
		nil,
		validMeasurementUnitConversionsBasePath,
		"from_unit",
		validMeasurementUnitID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidMeasurementUnitConversionsToUnitRequest builds an HTTP request for fetching a valid measurement conversion.
func (b *Builder) BuildGetValidMeasurementUnitConversionsToUnitRequest(ctx context.Context, validMeasurementUnitID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	uri := b.BuildURL(
		ctx,
		nil,
		validMeasurementUnitConversionsBasePath,
		"to_unit",
		validMeasurementUnitID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidMeasurementUnitConversionRequest builds an HTTP request for creating a valid measurement conversion.
func (b *Builder) BuildCreateValidMeasurementUnitConversionRequest(ctx context.Context, input *types.ValidMeasurementUnitConversionCreationRequestInput) (*http.Request, error) {
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
		validMeasurementUnitConversionsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidMeasurementUnitConversionRequest builds an HTTP request for updating a valid measurement conversion.
func (b *Builder) BuildUpdateValidMeasurementUnitConversionRequest(ctx context.Context, validMeasurementUnitConversion *types.ValidMeasurementUnitConversion) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementUnitConversion == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversion.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validMeasurementUnitConversionsBasePath,
		validMeasurementUnitConversion.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionUpdateRequestInput(validMeasurementUnitConversion)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidMeasurementUnitConversionRequest builds an HTTP request for archiving a valid measurement conversion.
func (b *Builder) BuildArchiveValidMeasurementUnitConversionRequest(ctx context.Context, validMeasurementUnitConversionID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementUnitConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	uri := b.BuildURL(
		ctx,
		nil,
		validMeasurementUnitConversionsBasePath,
		validMeasurementUnitConversionID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

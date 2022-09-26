package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	advancedPrepStepsBasePath = "advanced_prep_steps"
)

// BuildGetAdvancedPrepStepRequest builds an HTTP request for fetching a meal plan.
func (b *Builder) BuildGetAdvancedPrepStepRequest(ctx context.Context, advancedPrepStepID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if advancedPrepStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachAdvancedPrepStepIDToSpan(span, advancedPrepStepID)

	uri := b.BuildURL(
		ctx,
		nil,
		advancedPrepStepsBasePath,
		advancedPrepStepID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetAdvancedPrepStepsRequest builds an HTTP request for fetching a list of advanced prep steps.
func (b *Builder) BuildGetAdvancedPrepStepsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		advancedPrepStepsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildChangeAdvancedPrepStepStatusRequest builds an HTTP request for archiving a meal plan.
func (b *Builder) BuildChangeAdvancedPrepStepStatusRequest(ctx context.Context, input *types.AdvancedPrepStepStatusChangeRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachAdvancedPrepStepIDToSpan(span, input.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		advancedPrepStepsBasePath,
		input.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

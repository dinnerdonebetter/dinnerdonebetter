package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	webhooksBasePath = "webhooks"
)

// BuildGetWebhookRequest builds an HTTP request for fetching a webhook.
func (b *Builder) BuildGetWebhookRequest(ctx context.Context, webhookID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if webhookID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachWebhookIDToSpan(span, webhookID)

	uri := b.BuildURL(ctx, nil, webhooksBasePath, webhookID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetWebhooksRequest builds an HTTP request for fetching a list of webhooks.
func (b *Builder) BuildGetWebhooksRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)
	uri := b.BuildURL(ctx, filter.ToValues(), webhooksBasePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateWebhookRequest builds an HTTP request for creating a webhook.
func (b *Builder) BuildCreateWebhookRequest(ctx context.Context, input *types.WebhookCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(ctx, nil, webhooksBasePath)

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildArchiveWebhookRequest builds an HTTP request for archiving a webhook.
func (b *Builder) BuildArchiveWebhookRequest(ctx context.Context, webhookID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if webhookID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachWebhookIDToSpan(span, webhookID)

	uri := b.BuildURL(ctx, nil, webhooksBasePath, webhookID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

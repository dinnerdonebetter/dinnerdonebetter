package requests

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	serviceSettingsBasePath = "settings"
)

// BuildGetServiceSettingRequest builds an HTTP request for fetching a service setting.
func (b *Builder) BuildGetServiceSettingRequest(ctx context.Context, serviceSettingID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if serviceSettingID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachServiceSettingIDToSpan(span, serviceSettingID)

	uri := b.BuildURL(
		ctx,
		nil,
		serviceSettingsBasePath,
		serviceSettingID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildSearchServiceSettingsRequest builds an HTTP request for querying service settings.
func (b *Builder) BuildSearchServiceSettingsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger.WithValue(types.SearchQueryKey, query).WithValue(types.LimitQueryKey, limit)

	params := url.Values{}
	params.Set(types.SearchQueryKey, query)
	params.Set(types.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := b.BuildURL(
		ctx,
		params,
		serviceSettingsBasePath,
		"search",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		logger.Error(err, "building request")
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetServiceSettingsRequest builds an HTTP request for fetching a list of service settings.
func (b *Builder) BuildGetServiceSettingsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		serviceSettingsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		logger.Error(err, "building request")
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

package requests

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/panicking"
)

const (
	clientName = "ddb_client_v1"
)

// Builder is a client for interacting with v1 of our HTTP API.
type Builder struct {
	logger logging.Logger
	tracer tracing.Tracer
	url    *url.URL

	encoder  encoding.ClientEncoder
	panicker panicking.Panicker
}

// NewBuilder builds a new API client for us.
func NewBuilder(u *url.URL, logger logging.Logger, tracerProvider tracing.TracerProvider, encoder encoding.ClientEncoder) (*Builder, error) {
	l := logging.EnsureLogger(logger)

	if u == nil {
		return nil, ErrNoURLProvided
	}

	if encoder == nil {
		return nil, ErrNilEncoderProvided
	}

	c := &Builder{
		url:      u,
		logger:   l,
		encoder:  encoder,
		panicker: panicking.NewProductionPanicker(),
		tracer:   tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(clientName)),
	}

	return c, nil
}

// URL provides the client's URL.
func (b *Builder) URL() *url.URL {
	return b.url
}

// SetURL provides the client's URL.
func (b *Builder) SetURL(u *url.URL) error {
	if u == nil {
		return ErrNoURLProvided
	}

	b.url = u

	return nil
}

// BuildURL builds standard service URLs.
func (b *Builder) BuildURL(ctx context.Context, qp url.Values, parts ...string) string {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if u := b.buildAPIV1URL(ctx, qp, parts...); u != nil {
		return u.String()
	}

	return ""
}

// Must mandates that a given request is built without error.
func (b *Builder) Must(req *http.Request, err error) *http.Request {
	if err != nil {
		b.panicker.Panic(err)
	}

	return req
}

func buildRawURL(u *url.URL, qp url.Values, includeVersionPrefix bool, parts ...string) (*url.URL, error) {
	tu := *u

	if includeVersionPrefix {
		parts = append([]string{"api", "v1"}, parts...)
	}

	u, err := url.Parse(path.Join(parts...))
	if err != nil {
		return nil, observability.PrepareError(err, nil, "building URL")
	}

	if qp != nil {
		u.RawQuery = qp.Encode()
	}

	return tu.ResolveReference(u), nil
}

// buildRawURL takes a given set of query parameters and url parts, and returns a parsed url object from them.
func (b *Builder) buildAPIV1URL(ctx context.Context, queryParams url.Values, parts ...string) *url.URL {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tu := *b.url

	parts = append([]string{"api", "v1"}, parts...)

	u, err := url.Parse(path.Join(parts...))
	if err != nil {
		b.logger.Error(err, "building url")
		return nil
	}

	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	out := tu.ResolveReference(u)

	tracing.AttachToSpan(span, keys.RequestURIKey, out.String())

	return out
}

// buildUnversionedURL builds a url without the v1 API prefix. It should otherwise be identical to buildRawURL.
func (b *Builder) buildUnversionedURL(ctx context.Context, qp url.Values, parts ...string) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	u, err := buildRawURL(b.url, qp, false, parts...)
	if err != nil {
		b.logger.Error(err, "building unversioned url")
		return ""
	}

	return u.String()
}

// BuildWebsocketURL builds a standard url and then converts its scheme to the websocket protocol.
func (b *Builder) BuildWebsocketURL(ctx context.Context, parts ...string) string {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	u := b.buildAPIV1URL(ctx, nil, parts...)

	switch b.url.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	}

	return u.String()
}

// BuildHealthCheckRequest builds a health check HTTP request.
func (b *Builder) BuildHealthCheckRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	u := *b.url
	uri := fmt.Sprintf("%s://%s/_meta_/ready", u.Scheme, u.Host)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// buildDataRequest builds an HTTP request for a given method, url, and body data.
func (b *Builder) buildDataRequest(ctx context.Context, method, uri string, in any) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	body, err := b.encoder.EncodeReader(ctx, in)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encoding request")
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	req.Header.Set("RawHTML-type", b.encoder.ContentType())
	tracing.AttachToSpan(span, keys.RequestURIKey, req.URL.String())

	return req, nil
}

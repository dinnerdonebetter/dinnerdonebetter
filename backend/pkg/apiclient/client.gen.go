package apiclient

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"time"

	encoding "github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	logging "github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/moul/http2curl"
)

const (
	defaultTimeout = 30 * time.Second
	clientName     = "ddb_api_client_v1"

	zuckModeUserHeader      = "X-DDB-Zuck-Mode-User"
	zuckModeHouseholdHeader = "X-DDB-Zuck-Mode-Household"
)

// Client is a client for interacting with v1 of our HTTP API.
type Client struct {
	logger                  logging.Logger
	tracer                  tracing.Tracer
	contentType             encoding.ContentType
	encoder                 encoding.ClientEncoder
	authedClient            *http.Client
	unauthenticatedClient   *http.Client
	url                     *url.URL
	impersonatedUserID      string
	impersonatedHouseholdID string
	debug                   bool
}

// AuthenticatedClient returns the authenticated *apiclient.Client that we use to make most requests.
func (c *Client) AuthenticatedClient() *http.Client {
	return c.authedClient
}

// PlainClient returns the unauthenticated *apiclient.Client that we use to make certain requests.
func (c *Client) PlainClient() *http.Client {
	return c.unauthenticatedClient
}

// URL provides the client's URL.
func (c *Client) URL() *url.URL {
	return c.url
}

// NewClient builds a new API client for us.
func NewClient(u *url.URL, tracerProvider tracing.TracerProvider, options ...ClientOption) (*Client, error) {
	l := logging.NewNoopLogger()

	c := &Client{
		url:                   u,
		logger:                logging.EnsureLogger(nil),
		debug:                 false,
		contentType:           encoding.ContentTypeJSON,
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(clientName)),
		encoder:               encoding.ProvideClientEncoder(l, tracerProvider, encoding.ContentTypeJSON),
		authedClient:          tracing.BuildTracedHTTPClient(),
		unauthenticatedClient: tracing.BuildTracedHTTPClient(),
	}

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.url == nil || c.url.String() == "" {
		return nil, ErrNoURLProvided
	}

	return c, nil
}

// closeResponseBody takes a given HTTP response and closes its body, logging if an error occurs.
func (c *Client) closeResponseBody(ctx context.Context, res *http.Response) {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if res != nil {
		if err := res.Body.Close(); err != nil {
			observability.AcknowledgeError(err, c.logger.WithResponse(res), span, "closing response body")
		}
	}
}

// BuildURL builds standard service URLs.
func (c *Client) BuildURL(ctx context.Context, qp url.Values, parts ...string) string {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if u := c.buildRawURL(ctx, qp, parts...); u != nil {
		return u.String()
	}

	return ""
}

// buildRawURL takes a given set of query parameters and url parts, and returns a parsed url object from them.
func (c *Client) buildRawURL(ctx context.Context, queryParams url.Values, parts ...string) *url.URL {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	tu := *c.url
	logger := c.logger.WithValue(keys.URLQueryKey, queryParams.Encode())

	u, err := url.Parse(path.Join(parts...))
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "building URL")
		return nil
	}

	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	return tu.ResolveReference(u)
}

// IsUp returns whether the service's health endpoint is returning 200s.
func (c *Client) IsUp(ctx context.Context) bool {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BuildURL(ctx, nil, "_ops_", "ready"), http.NoBody)
	if err != nil {
		c.logger.Error("building steatus request", err)
		return false
	}

	res, err := c.fetchResponseToRequest(ctx, c.authedClient, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "performing health check")
		return false
	}
	defer c.closeResponseBody(ctx, res)

	return res.StatusCode == http.StatusOK
}

func (c *Client) logRequest(logger logging.Logger, res *http.Response) {
	if bdump, err := httputil.DumpResponse(res, true); err == nil {
		logger = logger.WithValue("response_body", string(bdump))
	}

	logger.WithValue(keys.ResponseStatusKey, res.StatusCode).Debug("request executed")
}

// fetchResponseToRequest takes a given *http.Request and executes it with the provided.
// client, alongside some debugging logging.
func (c *Client) fetchResponseToRequest(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.WithRequest(req).WithSpan(span)

	if command, err := http2curl.GetCurlCommand(req); err == nil && c.debug {
		logger = logger.WithValue("curl", command.String())
	}

	// this should be the only use of .Do in this package
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing request")
	}

	return res, nil
}

// executeAndUnmarshal executes a request and unmarshals it to the provided interface.
func (c *Client) executeAndUnmarshal(ctx context.Context, req *http.Request, httpClient *http.Client, out any) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.WithRequest(req).WithSpan(span)
	logger.Debug("executing request")

	res, err := c.fetchResponseToRequest(ctx, httpClient, req)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "executing request")
	}

	if resErr := errorFromResponse(res); resErr != nil {
		return resErr
	}

	if out != nil {
		if err = c.unmarshalBody(ctx, res, out); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "%s %s %d", res.Request.Method, res.Request.URL.Path, res.StatusCode)
		}
	}

	if err = res.Body.Close(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "closing response body")
	}

	return nil
}

// fetchAndUnmarshal takes a given request and executes it with the auth client.
func (c *Client) fetchAndUnmarshal(ctx context.Context, req *http.Request, out any) error {
	return c.executeAndUnmarshal(ctx, req, c.authedClient, out)
}

// fetchAndUnmarshalWithoutAuthentication takes a given request and executes it with the plain client.
func (c *Client) fetchAndUnmarshalWithoutAuthentication(ctx context.Context, req *http.Request, out any) error {
	return c.executeAndUnmarshal(ctx, req, c.unauthenticatedClient, out)
}

// responseIsOK executes an HTTP request and loads the response content into a bool.
func (c *Client) responseIsOK(ctx context.Context, req *http.Request) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	res, err := c.fetchResponseToRequest(ctx, c.authedClient, req)
	if err != nil {
		return false, observability.PrepareError(err, span, "executing existence request")
	}

	c.closeResponseBody(ctx, res)

	return res.StatusCode == http.StatusOK, nil
}

// buildDataRequest builds an HTTP request for a given method, url, and body data.
func (c *Client) buildDataRequest(ctx context.Context, method, uri string, in any) (*http.Request, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	body, err := c.encoder.EncodeReader(ctx, in)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encoding request")
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	tracing.AttachToSpan(span, keys.RequestURIKey, req.URL.String())

	return req, nil
}

type RequestModifier func(*http.Request)

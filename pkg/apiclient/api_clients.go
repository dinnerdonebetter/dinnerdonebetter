package apiclient

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetAPIClient gets an API client.
func (c *Client) GetAPIClient(ctx context.Context, apiClientDatabaseID string) (*types.APIClient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if apiClientDatabaseID == "" {
		return nil, ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildGetAPIClientRequest(ctx, apiClientDatabaseID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building retrieve API client request")
	}

	var apiClient *types.APIClient
	if err = c.fetchAndUnmarshal(ctx, req, &apiClient); err != nil {
		return nil, observability.PrepareError(err, span, "fetching api client")
	}

	return apiClient, nil
}

// GetAPIClients gets a list of API clients.
func (c *Client) GetAPIClients(ctx context.Context, filter *types.QueryFilter) (*types.APIClientList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetAPIClientsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building API clients list request")
	}

	var apiClients *types.APIClientList
	if err = c.fetchAndUnmarshal(ctx, req, &apiClients); err != nil {
		return nil, observability.PrepareError(err, span, "fetching api clients")
	}

	return apiClients, nil
}

// CreateAPIClient creates an API client.
func (c *Client) CreateAPIClient(ctx context.Context, cookie *http.Cookie, input *types.APIClientCreationRequestInput) (*types.APIClientCreationResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if cookie == nil && c.authMethod != cookieAuthMethod {
		return nil, ErrCookieRequired
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	var apiClientResponse *types.APIClientCreationResponse

	req, err := c.requestBuilder.BuildCreateAPIClientRequest(ctx, cookie, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building create API client request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &apiClientResponse); err != nil {
		return nil, observability.PrepareError(err, span, "creating api client")
	}

	return apiClientResponse, nil
}

// ArchiveAPIClient archives an API client.
func (c *Client) ArchiveAPIClient(ctx context.Context, apiClientDatabaseID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if apiClientDatabaseID == "" {
		return ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildArchiveAPIClientRequest(ctx, apiClientDatabaseID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive API client request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "archiving api client")
	}

	return nil
}

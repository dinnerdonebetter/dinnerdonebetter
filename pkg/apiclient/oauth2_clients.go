package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetOAuth2Client gets an OAuth2 client.
func (c *Client) GetOAuth2Client(ctx context.Context, oauth2ClientDatabaseID string) (*types.OAuth2Client, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if oauth2ClientDatabaseID == "" {
		return nil, ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildGetOAuth2ClientRequest(ctx, oauth2ClientDatabaseID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building retrieve OAuth2 client request")
	}

	var oauth2Client *types.OAuth2Client
	if err = c.fetchAndUnmarshal(ctx, req, &oauth2Client); err != nil {
		return nil, observability.PrepareError(err, span, "fetching api client")
	}

	return oauth2Client, nil
}

// GetOAuth2Clients gets a list of OAuth2 clients.
func (c *Client) GetOAuth2Clients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.OAuth2Client], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetOAuth2ClientsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building OAuth2 clients list request")
	}

	var oauth2Clients *types.QueryFilteredResult[types.OAuth2Client]
	if err = c.fetchAndUnmarshal(ctx, req, &oauth2Clients); err != nil {
		return nil, observability.PrepareError(err, span, "fetching api clients")
	}

	return oauth2Clients, nil
}

// CreateOAuth2Client creates an OAuth2 client.
func (c *Client) CreateOAuth2Client(ctx context.Context, input *types.OAuth2ClientCreationRequestInput) (*types.OAuth2ClientCreationResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if c.authMethod != cookieAuthMethod {
		return nil, ErrCookieRequired
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	var oauth2ClientResponse *types.OAuth2ClientCreationResponse

	req, err := c.requestBuilder.BuildCreateOAuth2ClientRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building create OAuth2 client request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &oauth2ClientResponse); err != nil {
		return nil, observability.PrepareError(err, span, "creating api client")
	}

	return oauth2ClientResponse, nil
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (c *Client) ArchiveOAuth2Client(ctx context.Context, oauth2ClientDatabaseID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if oauth2ClientDatabaseID == "" {
		return ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildArchiveOAuth2ClientRequest(ctx, oauth2ClientDatabaseID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive OAuth2 client request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "archiving api client")
	}

	return nil
}

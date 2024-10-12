package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetOAuth2Client gets an OAuth2 client.
func (c *Client) GetOAuth2Client(ctx context.Context, oauth2ClientDatabaseID string) (*types.OAuth2Client, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if oauth2ClientDatabaseID == "" {
		return nil, ErrInvalidIDProvided
	}

	res, err := c.authedGeneratedClient.GetOAuth2Client(ctx, oauth2ClientDatabaseID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "retrieve OAuth2 client")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.OAuth2Client]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "fetching oauth2 client")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetOAuth2Clients gets a list of OAuth2 clients.
func (c *Client) GetOAuth2Clients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.OAuth2Client], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetOAuth2ClientsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetOAuth2Clients(ctx, params)
	if err != nil {
		return nil, observability.PrepareError(err, span, "OAuth2 clients list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.OAuth2Client]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "fetching api clients")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.OAuth2Client]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateOAuth2Client creates an OAuth2 client.
func (c *Client) CreateOAuth2Client(ctx context.Context, input *types.OAuth2ClientCreationRequestInput) (*types.OAuth2ClientCreationResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	body := generated.CreateOAuth2ClientJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateOAuth2Client(ctx, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "create OAuth2 client")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.OAuth2ClientCreationResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "creating oauth2 client")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (c *Client) ArchiveOAuth2Client(ctx context.Context, oauth2ClientDatabaseID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if oauth2ClientDatabaseID == "" {
		return ErrInvalidIDProvided
	}

	res, err := c.authedGeneratedClient.ArchiveOAuth2Client(ctx, oauth2ClientDatabaseID)
	if err != nil {
		return observability.PrepareError(err, span, "archive OAuth2 client")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.OAuth2Client]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "archiving oauth2 client")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

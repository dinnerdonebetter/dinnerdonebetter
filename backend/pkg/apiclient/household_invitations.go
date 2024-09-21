package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetHouseholdInvitation retrieves a household invitation.
func (c *Client) GetHouseholdInvitation(ctx context.Context, householdInvitationID string) (*types.HouseholdInvitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInvitationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	res, err := c.authedGeneratedClient.GetHouseholdInvitation(ctx, householdInvitationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get invitation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetPendingHouseholdInvitationsFromUser retrieves household invitations sent by the user.
func (c *Client) GetPendingHouseholdInvitationsFromUser(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInvitation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	filter.AttachToLogger(logger)

	params := &generated.GetSentHouseholdInvitationsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetSentHouseholdInvitations(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "reject invitation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.HouseholdInvitation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.HouseholdInvitation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetPendingHouseholdInvitationsForUser retrieves household invitations received by the user.
func (c *Client) GetPendingHouseholdInvitationsForUser(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInvitation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	filter.AttachToLogger(logger)

	params := &generated.GetReceivedHouseholdInvitationsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetReceivedHouseholdInvitations(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "reject invitation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.HouseholdInvitation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.HouseholdInvitation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// AcceptHouseholdInvitation accepts a given household invitation.
func (c *Client) AcceptHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	if token == "" {
		return ErrInvalidIDProvided
	}

	body := generated.AcceptHouseholdInvitationJSONRequestBody{
		Note:  &note,
		Token: &token,
	}

	res, err := c.authedGeneratedClient.AcceptHouseholdInvitation(ctx, householdInvitationID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "accept invitation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// CancelHouseholdInvitation cancels a given household invitation.
func (c *Client) CancelHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	if token == "" {
		return ErrInvalidIDProvided
	}

	body := generated.CancelHouseholdInvitationJSONRequestBody{
		Note:  &note,
		Token: &token,
	}

	res, err := c.authedGeneratedClient.CancelHouseholdInvitation(ctx, householdInvitationID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "cancel invitation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// RejectHouseholdInvitation rejects a given household invitation.
func (c *Client) RejectHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	if token == "" {
		return ErrInvalidIDProvided
	}

	body := generated.RejectHouseholdInvitationJSONRequestBody{
		Note:  &note,
		Token: &token,
	}

	res, err := c.authedGeneratedClient.RejectHouseholdInvitation(ctx, householdInvitationID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "reject invitation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

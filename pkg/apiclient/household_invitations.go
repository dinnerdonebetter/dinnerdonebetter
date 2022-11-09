package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/pkg/types"
)

// GetHouseholdInvitation retrieves a household invitation.
func (c *Client) GetHouseholdInvitation(ctx context.Context, householdID, householdInvitationID string) (*types.HouseholdInvitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	if householdInvitationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	req, err := c.requestBuilder.BuildGetHouseholdInvitationRequest(ctx, householdID, householdInvitationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get invitation request")
	}

	var invitation *types.HouseholdInvitation
	if err = c.fetchAndUnmarshal(ctx, req, &invitation); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving invitation")
	}

	return invitation, nil
}

// GetPendingHouseholdInvitationsFromUser retrieves household invitations sent by the user.
func (c *Client) GetPendingHouseholdInvitationsFromUser(ctx context.Context, filter *types.QueryFilter) (*types.HouseholdInvitationList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	filter.AttachToLogger(logger)

	req, err := c.requestBuilder.BuildGetPendingHouseholdInvitationsFromUserRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	var invitationList *types.HouseholdInvitationList
	if err = c.fetchAndUnmarshal(ctx, req, &invitationList); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	return invitationList, nil
}

// GetPendingHouseholdInvitationsForUser retrieves household invitations received by the user.
func (c *Client) GetPendingHouseholdInvitationsForUser(ctx context.Context, filter *types.QueryFilter) (*types.HouseholdInvitationList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	filter.AttachToLogger(logger)

	req, err := c.requestBuilder.BuildGetPendingHouseholdInvitationsForUserRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	var invitationList *types.HouseholdInvitationList
	if err = c.fetchAndUnmarshal(ctx, req, &invitationList); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	return invitationList, nil
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

	req, err := c.requestBuilder.BuildAcceptHouseholdInvitationRequest(ctx, householdInvitationID, token, note)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
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

	req, err := c.requestBuilder.BuildCancelHouseholdInvitationRequest(ctx, householdInvitationID, token, note)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
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

	req, err := c.requestBuilder.BuildRejectHouseholdInvitationRequest(ctx, householdInvitationID, token, note)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	return nil
}

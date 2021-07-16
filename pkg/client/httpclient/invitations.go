package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// InvitationExists retrieves whether an invitation exists.
func (c *Client) InvitationExists(ctx context.Context, invitationID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if invitationID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	req, err := c.requestBuilder.BuildInvitationExistsRequest(ctx, invitationID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building invitation existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for invitation #%d", invitationID)
	}

	return exists, nil
}

// GetInvitation gets an invitation.
func (c *Client) GetInvitation(ctx context.Context, invitationID uint64) (*types.Invitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if invitationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	req, err := c.requestBuilder.BuildGetInvitationRequest(ctx, invitationID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get invitation request")
	}

	var invitation *types.Invitation
	if err = c.fetchAndUnmarshal(ctx, req, &invitation); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving invitation")
	}

	return invitation, nil
}

// GetInvitations retrieves a list of invitations.
func (c *Client) GetInvitations(ctx context.Context, filter *types.QueryFilter) (*types.InvitationList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetInvitationsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building invitations list request")
	}

	var invitations *types.InvitationList
	if err = c.fetchAndUnmarshal(ctx, req, &invitations); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving invitations")
	}

	return invitations, nil
}

// CreateInvitation creates an invitation.
func (c *Client) CreateInvitation(ctx context.Context, input *types.InvitationCreationInput) (*types.Invitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateInvitationRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create invitation request")
	}

	var invitation *types.Invitation
	if err = c.fetchAndUnmarshal(ctx, req, &invitation); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating invitation")
	}

	return invitation, nil
}

// UpdateInvitation updates an invitation.
func (c *Client) UpdateInvitation(ctx context.Context, invitation *types.Invitation) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if invitation == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitation.ID)
	tracing.AttachInvitationIDToSpan(span, invitation.ID)

	req, err := c.requestBuilder.BuildUpdateInvitationRequest(ctx, invitation)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update invitation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &invitation); err != nil {
		return observability.PrepareError(err, logger, span, "updating invitation #%d", invitation.ID)
	}

	return nil
}

// ArchiveInvitation archives an invitation.
func (c *Client) ArchiveInvitation(ctx context.Context, invitationID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if invitationID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	req, err := c.requestBuilder.BuildArchiveInvitationRequest(ctx, invitationID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive invitation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving invitation #%d", invitationID)
	}

	return nil
}

// GetAuditLogForInvitation retrieves a list of audit log entries pertaining to an invitation.
func (c *Client) GetAuditLogForInvitation(ctx context.Context, invitationID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if invitationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	req, err := c.requestBuilder.BuildGetAuditLogForInvitationRequest(ctx, invitationID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for invitation request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

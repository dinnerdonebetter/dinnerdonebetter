package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	invitationsBasePath = "invitations"
)

// BuildInvitationExistsRequest builds an HTTP request for checking the existence of an invitation.
func (c *V1Client) BuildInvitationExistsRequest(ctx context.Context, invitationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildInvitationExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		invitationsBasePath,
		strconv.FormatUint(invitationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// InvitationExists retrieves whether or not an invitation exists.
func (c *V1Client) InvitationExists(ctx context.Context, invitationID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "InvitationExists")
	defer span.End()

	req, err := c.BuildInvitationExistsRequest(ctx, invitationID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetInvitationRequest builds an HTTP request for fetching an invitation.
func (c *V1Client) BuildGetInvitationRequest(ctx context.Context, invitationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetInvitationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		invitationsBasePath,
		strconv.FormatUint(invitationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetInvitation retrieves an invitation.
func (c *V1Client) GetInvitation(ctx context.Context, invitationID uint64) (invitation *models.Invitation, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetInvitation")
	defer span.End()

	req, err := c.BuildGetInvitationRequest(ctx, invitationID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &invitation); retrieveErr != nil {
		return nil, retrieveErr
	}

	return invitation, nil
}

// BuildGetInvitationsRequest builds an HTTP request for fetching invitations.
func (c *V1Client) BuildGetInvitationsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetInvitationsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		invitationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetInvitations retrieves a list of invitations.
func (c *V1Client) GetInvitations(ctx context.Context, filter *models.QueryFilter) (invitations *models.InvitationList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetInvitations")
	defer span.End()

	req, err := c.BuildGetInvitationsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &invitations); retrieveErr != nil {
		return nil, retrieveErr
	}

	return invitations, nil
}

// BuildCreateInvitationRequest builds an HTTP request for creating an invitation.
func (c *V1Client) BuildCreateInvitationRequest(ctx context.Context, input *models.InvitationCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateInvitationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		invitationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateInvitation creates an invitation.
func (c *V1Client) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (invitation *models.Invitation, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateInvitation")
	defer span.End()

	req, err := c.BuildCreateInvitationRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &invitation)
	return invitation, err
}

// BuildUpdateInvitationRequest builds an HTTP request for updating an invitation.
func (c *V1Client) BuildUpdateInvitationRequest(ctx context.Context, invitation *models.Invitation) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateInvitationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		invitationsBasePath,
		strconv.FormatUint(invitation.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, invitation)
}

// UpdateInvitation updates an invitation.
func (c *V1Client) UpdateInvitation(ctx context.Context, invitation *models.Invitation) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateInvitation")
	defer span.End()

	req, err := c.BuildUpdateInvitationRequest(ctx, invitation)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &invitation)
}

// BuildArchiveInvitationRequest builds an HTTP request for updating an invitation.
func (c *V1Client) BuildArchiveInvitationRequest(ctx context.Context, invitationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveInvitationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		invitationsBasePath,
		strconv.FormatUint(invitationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveInvitation archives an invitation.
func (c *V1Client) ArchiveInvitation(ctx context.Context, invitationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveInvitation")
	defer span.End()

	req, err := c.BuildArchiveInvitationRequest(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

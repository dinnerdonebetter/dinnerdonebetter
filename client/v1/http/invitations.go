package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	invitationsBasePath = "invitations"
)

// BuildGetInvitationRequest builds an HTTP request for fetching an invitation
func (c *V1Client) BuildGetInvitationRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, invitationsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetInvitation retrieves an invitation
func (c *V1Client) GetInvitation(ctx context.Context, id uint64) (invitation *models.Invitation, err error) {
	req, err := c.BuildGetInvitationRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &invitation); retrieveErr != nil {
		return nil, retrieveErr
	}

	return invitation, nil
}

// BuildGetInvitationsRequest builds an HTTP request for fetching invitations
func (c *V1Client) BuildGetInvitationsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), invitationsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetInvitations retrieves a list of invitations
func (c *V1Client) GetInvitations(ctx context.Context, filter *models.QueryFilter) (invitations *models.InvitationList, err error) {
	req, err := c.BuildGetInvitationsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &invitations); retrieveErr != nil {
		return nil, retrieveErr
	}

	return invitations, nil
}

// BuildCreateInvitationRequest builds an HTTP request for creating an invitation
func (c *V1Client) BuildCreateInvitationRequest(ctx context.Context, body *models.InvitationCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, invitationsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateInvitation creates an invitation
func (c *V1Client) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (invitation *models.Invitation, err error) {
	req, err := c.BuildCreateInvitationRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &invitation)
	return invitation, err
}

// BuildUpdateInvitationRequest builds an HTTP request for updating an invitation
func (c *V1Client) BuildUpdateInvitationRequest(ctx context.Context, updated *models.Invitation) (*http.Request, error) {
	uri := c.BuildURL(nil, invitationsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateInvitation updates an invitation
func (c *V1Client) UpdateInvitation(ctx context.Context, updated *models.Invitation) error {
	req, err := c.BuildUpdateInvitationRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveInvitationRequest builds an HTTP request for updating an invitation
func (c *V1Client) BuildArchiveInvitationRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, invitationsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveInvitation archives an invitation
func (c *V1Client) ArchiveInvitation(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveInvitationRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

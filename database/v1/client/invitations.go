package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.InvitationDataManager = (*Client)(nil)

// attachInvitationIDToSpan provides a consistent way to attach an invitation's ID to a span
func attachInvitationIDToSpan(span *trace.Span, invitationID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("invitation_id", strconv.FormatUint(invitationID, 10)))
	}
}

// GetInvitation fetches an invitation from the database
func (c *Client) GetInvitation(ctx context.Context, invitationID, userID uint64) (*models.Invitation, error) {
	ctx, span := trace.StartSpan(ctx, "GetInvitation")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachInvitationIDToSpan(span, invitationID)

	c.logger.WithValues(map[string]interface{}{
		"invitation_id": invitationID,
		"user_id":       userID,
	}).Debug("GetInvitation called")

	return c.querier.GetInvitation(ctx, invitationID, userID)
}

// GetInvitationCount fetches the count of invitations from the database that meet a particular filter
func (c *Client) GetInvitationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetInvitationCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetInvitationCount called")

	return c.querier.GetInvitationCount(ctx, filter, userID)
}

// GetAllInvitationsCount fetches the count of invitations from the database that meet a particular filter
func (c *Client) GetAllInvitationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllInvitationsCount")
	defer span.End()

	c.logger.Debug("GetAllInvitationsCount called")

	return c.querier.GetAllInvitationsCount(ctx)
}

// GetInvitations fetches a list of invitations from the database that meet a particular filter
func (c *Client) GetInvitations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.InvitationList, error) {
	ctx, span := trace.StartSpan(ctx, "GetInvitations")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetInvitations called")

	invitationList, err := c.querier.GetInvitations(ctx, filter, userID)

	return invitationList, err
}

// GetAllInvitationsForUser fetches a list of invitations from the database that meet a particular filter
func (c *Client) GetAllInvitationsForUser(ctx context.Context, userID uint64) ([]models.Invitation, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllInvitationsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllInvitationsForUser called")

	invitationList, err := c.querier.GetAllInvitationsForUser(ctx, userID)

	return invitationList, err
}

// CreateInvitation creates an invitation in the database
func (c *Client) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (*models.Invitation, error) {
	ctx, span := trace.StartSpan(ctx, "CreateInvitation")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateInvitation called")

	return c.querier.CreateInvitation(ctx, input)
}

// UpdateInvitation updates a particular invitation. Note that UpdateInvitation expects the
// provided input to have a valid ID.
func (c *Client) UpdateInvitation(ctx context.Context, input *models.Invitation) error {
	ctx, span := trace.StartSpan(ctx, "UpdateInvitation")
	defer span.End()

	attachInvitationIDToSpan(span, input.ID)
	c.logger.WithValue("invitation_id", input.ID).Debug("UpdateInvitation called")

	return c.querier.UpdateInvitation(ctx, input)
}

// ArchiveInvitation archives an invitation from the database by its ID
func (c *Client) ArchiveInvitation(ctx context.Context, invitationID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveInvitation")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachInvitationIDToSpan(span, invitationID)

	c.logger.WithValues(map[string]interface{}{
		"invitation_id": invitationID,
		"user_id":       userID,
	}).Debug("ArchiveInvitation called")

	return c.querier.ArchiveInvitation(ctx, invitationID, userID)
}

package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.InvitationDataManager = (*Client)(nil)

// InvitationExists fetches whether or not an invitation exists from the database.
func (c *Client) InvitationExists(ctx context.Context, invitationID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "InvitationExists")
	defer span.End()

	tracing.AttachInvitationIDToSpan(span, invitationID)

	c.logger.WithValues(map[string]interface{}{
		"invitation_id": invitationID,
	}).Debug("InvitationExists called")

	return c.querier.InvitationExists(ctx, invitationID)
}

// GetInvitation fetches an invitation from the database.
func (c *Client) GetInvitation(ctx context.Context, invitationID uint64) (*models.Invitation, error) {
	ctx, span := tracing.StartSpan(ctx, "GetInvitation")
	defer span.End()

	tracing.AttachInvitationIDToSpan(span, invitationID)

	c.logger.WithValues(map[string]interface{}{
		"invitation_id": invitationID,
	}).Debug("GetInvitation called")

	return c.querier.GetInvitation(ctx, invitationID)
}

// GetAllInvitationsCount fetches the count of invitations from the database that meet a particular filter.
func (c *Client) GetAllInvitationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllInvitationsCount")
	defer span.End()

	c.logger.Debug("GetAllInvitationsCount called")

	return c.querier.GetAllInvitationsCount(ctx)
}

// GetAllInvitations fetches a list of all invitations in the database.
func (c *Client) GetAllInvitations(ctx context.Context, results chan []models.Invitation) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllInvitations")
	defer span.End()

	c.logger.Debug("GetAllInvitations called")

	return c.querier.GetAllInvitations(ctx, results)
}

// GetInvitations fetches a list of invitations from the database that meet a particular filter.
func (c *Client) GetInvitations(ctx context.Context, filter *models.QueryFilter) (*models.InvitationList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetInvitations")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetInvitations called")

	invitationList, err := c.querier.GetInvitations(ctx, filter)

	return invitationList, err
}

// GetInvitationsWithIDs fetches invitations from the database within a given set of IDs.
func (c *Client) GetInvitationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.Invitation, error) {
	ctx, span := tracing.StartSpan(ctx, "GetInvitationsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetInvitationsWithIDs called")

	invitationList, err := c.querier.GetInvitationsWithIDs(ctx, limit, ids)

	return invitationList, err
}

// CreateInvitation creates an invitation in the database.
func (c *Client) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (*models.Invitation, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateInvitation")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateInvitation called")

	return c.querier.CreateInvitation(ctx, input)
}

// UpdateInvitation updates a particular invitation. Note that UpdateInvitation expects the
// provided input to have a valid ID.
func (c *Client) UpdateInvitation(ctx context.Context, updated *models.Invitation) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateInvitation")
	defer span.End()

	tracing.AttachInvitationIDToSpan(span, updated.ID)
	c.logger.WithValue("invitation_id", updated.ID).Debug("UpdateInvitation called")

	return c.querier.UpdateInvitation(ctx, updated)
}

// ArchiveInvitation archives an invitation from the database by its ID.
func (c *Client) ArchiveInvitation(ctx context.Context, invitationID, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveInvitation")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	c.logger.WithValues(map[string]interface{}{
		"invitation_id": invitationID,
		"user_id":       userID,
	}).Debug("ArchiveInvitation called")

	return c.querier.ArchiveInvitation(ctx, invitationID, userID)
}

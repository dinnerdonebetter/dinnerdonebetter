package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.PreparationDataManager = (*Client)(nil)

// attachPreparationIDToSpan provides a consistent way to attach a preparation's ID to a span
func attachPreparationIDToSpan(span *trace.Span, preparationID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("preparation_id", strconv.FormatUint(preparationID, 10)))
	}
}

// GetPreparation fetches a preparation from the database
func (c *Client) GetPreparation(ctx context.Context, preparationID, userID uint64) (*models.Preparation, error) {
	ctx, span := trace.StartSpan(ctx, "GetPreparation")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachPreparationIDToSpan(span, preparationID)

	c.logger.WithValues(map[string]interface{}{
		"preparation_id": preparationID,
		"user_id":        userID,
	}).Debug("GetPreparation called")

	return c.querier.GetPreparation(ctx, preparationID, userID)
}

// GetPreparationCount fetches the count of preparations from the database that meet a particular filter
func (c *Client) GetPreparationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetPreparationCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetPreparationCount called")

	return c.querier.GetPreparationCount(ctx, filter, userID)
}

// GetAllPreparationsCount fetches the count of preparations from the database that meet a particular filter
func (c *Client) GetAllPreparationsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllPreparationsCount")
	defer span.End()

	c.logger.Debug("GetAllPreparationsCount called")

	return c.querier.GetAllPreparationsCount(ctx)
}

// GetPreparations fetches a list of preparations from the database that meet a particular filter
func (c *Client) GetPreparations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.PreparationList, error) {
	ctx, span := trace.StartSpan(ctx, "GetPreparations")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetPreparations called")

	preparationList, err := c.querier.GetPreparations(ctx, filter, userID)

	return preparationList, err
}

// GetAllPreparationsForUser fetches a list of preparations from the database that meet a particular filter
func (c *Client) GetAllPreparationsForUser(ctx context.Context, userID uint64) ([]models.Preparation, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllPreparationsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllPreparationsForUser called")

	preparationList, err := c.querier.GetAllPreparationsForUser(ctx, userID)

	return preparationList, err
}

// CreatePreparation creates a preparation in the database
func (c *Client) CreatePreparation(ctx context.Context, input *models.PreparationCreationInput) (*models.Preparation, error) {
	ctx, span := trace.StartSpan(ctx, "CreatePreparation")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreatePreparation called")

	return c.querier.CreatePreparation(ctx, input)
}

// UpdatePreparation updates a particular preparation. Note that UpdatePreparation expects the
// provided input to have a valid ID.
func (c *Client) UpdatePreparation(ctx context.Context, input *models.Preparation) error {
	ctx, span := trace.StartSpan(ctx, "UpdatePreparation")
	defer span.End()

	attachPreparationIDToSpan(span, input.ID)
	c.logger.WithValue("preparation_id", input.ID).Debug("UpdatePreparation called")

	return c.querier.UpdatePreparation(ctx, input)
}

// ArchivePreparation archives a preparation from the database by its ID
func (c *Client) ArchivePreparation(ctx context.Context, preparationID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchivePreparation")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachPreparationIDToSpan(span, preparationID)

	c.logger.WithValues(map[string]interface{}{
		"preparation_id": preparationID,
		"user_id":        userID,
	}).Debug("ArchivePreparation called")

	return c.querier.ArchivePreparation(ctx, preparationID, userID)
}

package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.IterationMediaDataManager = (*Client)(nil)

// attachIterationMediaIDToSpan provides a consistent way to attach an iteration media's ID to a span
func attachIterationMediaIDToSpan(span *trace.Span, iterationMediaID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("iteration_media_id", strconv.FormatUint(iterationMediaID, 10)))
	}
}

// GetIterationMedia fetches an iteration media from the database
func (c *Client) GetIterationMedia(ctx context.Context, iterationMediaID, userID uint64) (*models.IterationMedia, error) {
	ctx, span := trace.StartSpan(ctx, "GetIterationMedia")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachIterationMediaIDToSpan(span, iterationMediaID)

	c.logger.WithValues(map[string]interface{}{
		"iteration_media_id": iterationMediaID,
		"user_id":            userID,
	}).Debug("GetIterationMedia called")

	return c.querier.GetIterationMedia(ctx, iterationMediaID, userID)
}

// GetIterationMediaCount fetches the count of iteration medias from the database that meet a particular filter
func (c *Client) GetIterationMediaCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetIterationMediaCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetIterationMediaCount called")

	return c.querier.GetIterationMediaCount(ctx, filter, userID)
}

// GetAllIterationMediasCount fetches the count of iteration medias from the database that meet a particular filter
func (c *Client) GetAllIterationMediasCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllIterationMediasCount")
	defer span.End()

	c.logger.Debug("GetAllIterationMediasCount called")

	return c.querier.GetAllIterationMediasCount(ctx)
}

// GetIterationMedias fetches a list of iteration medias from the database that meet a particular filter
func (c *Client) GetIterationMedias(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.IterationMediaList, error) {
	ctx, span := trace.StartSpan(ctx, "GetIterationMedias")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetIterationMedias called")

	iterationMediaList, err := c.querier.GetIterationMedias(ctx, filter, userID)

	return iterationMediaList, err
}

// GetAllIterationMediasForUser fetches a list of iteration medias from the database that meet a particular filter
func (c *Client) GetAllIterationMediasForUser(ctx context.Context, userID uint64) ([]models.IterationMedia, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllIterationMediasForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllIterationMediasForUser called")

	iterationMediaList, err := c.querier.GetAllIterationMediasForUser(ctx, userID)

	return iterationMediaList, err
}

// CreateIterationMedia creates an iteration media in the database
func (c *Client) CreateIterationMedia(ctx context.Context, input *models.IterationMediaCreationInput) (*models.IterationMedia, error) {
	ctx, span := trace.StartSpan(ctx, "CreateIterationMedia")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateIterationMedia called")

	return c.querier.CreateIterationMedia(ctx, input)
}

// UpdateIterationMedia updates a particular iteration media. Note that UpdateIterationMedia expects the
// provided input to have a valid ID.
func (c *Client) UpdateIterationMedia(ctx context.Context, input *models.IterationMedia) error {
	ctx, span := trace.StartSpan(ctx, "UpdateIterationMedia")
	defer span.End()

	attachIterationMediaIDToSpan(span, input.ID)
	c.logger.WithValue("iteration_media_id", input.ID).Debug("UpdateIterationMedia called")

	return c.querier.UpdateIterationMedia(ctx, input)
}

// ArchiveIterationMedia archives an iteration media from the database by its ID
func (c *Client) ArchiveIterationMedia(ctx context.Context, iterationMediaID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveIterationMedia")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachIterationMediaIDToSpan(span, iterationMediaID)

	c.logger.WithValues(map[string]interface{}{
		"iteration_media_id": iterationMediaID,
		"user_id":            userID,
	}).Debug("ArchiveIterationMedia called")

	return c.querier.ArchiveIterationMedia(ctx, iterationMediaID, userID)
}

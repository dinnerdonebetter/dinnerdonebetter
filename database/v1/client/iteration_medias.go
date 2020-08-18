package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.IterationMediaDataManager = (*Client)(nil)

// IterationMediaExists fetches whether or not an iteration media exists from the database.
func (c *Client) IterationMediaExists(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "IterationMediaExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
	tracing.AttachIterationMediaIDToSpan(span, iterationMediaID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":           recipeID,
		"recipe_iteration_id": recipeIterationID,
		"iteration_media_id":  iterationMediaID,
	}).Debug("IterationMediaExists called")

	return c.querier.IterationMediaExists(ctx, recipeID, recipeIterationID, iterationMediaID)
}

// GetIterationMedia fetches an iteration media from the database.
func (c *Client) GetIterationMedia(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (*models.IterationMedia, error) {
	ctx, span := tracing.StartSpan(ctx, "GetIterationMedia")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
	tracing.AttachIterationMediaIDToSpan(span, iterationMediaID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":           recipeID,
		"recipe_iteration_id": recipeIterationID,
		"iteration_media_id":  iterationMediaID,
	}).Debug("GetIterationMedia called")

	return c.querier.GetIterationMedia(ctx, recipeID, recipeIterationID, iterationMediaID)
}

// GetAllIterationMediasCount fetches the count of iteration medias from the database that meet a particular filter.
func (c *Client) GetAllIterationMediasCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllIterationMediasCount")
	defer span.End()

	c.logger.Debug("GetAllIterationMediasCount called")

	return c.querier.GetAllIterationMediasCount(ctx)
}

// GetAllIterationMedias fetches a list of all iteration medias in the database.
func (c *Client) GetAllIterationMedias(ctx context.Context, results chan []models.IterationMedia) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllIterationMedias")
	defer span.End()

	c.logger.Debug("GetAllIterationMedias called")

	return c.querier.GetAllIterationMedias(ctx, results)
}

// GetIterationMedias fetches a list of iteration medias from the database that meet a particular filter.
func (c *Client) GetIterationMedias(ctx context.Context, recipeID, recipeIterationID uint64, filter *models.QueryFilter) (*models.IterationMediaList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetIterationMedias")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":           recipeID,
		"recipe_iteration_id": recipeIterationID,
	}).Debug("GetIterationMedias called")

	iterationMediaList, err := c.querier.GetIterationMedias(ctx, recipeID, recipeIterationID, filter)

	return iterationMediaList, err
}

// GetIterationMediasWithIDs fetches iteration medias from the database within a given set of IDs.
func (c *Client) GetIterationMediasWithIDs(ctx context.Context, recipeID, recipeIterationID uint64, limit uint8, ids []uint64) ([]models.IterationMedia, error) {
	ctx, span := tracing.StartSpan(ctx, "GetIterationMediasWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetIterationMediasWithIDs called")

	iterationMediaList, err := c.querier.GetIterationMediasWithIDs(ctx, recipeID, recipeIterationID, limit, ids)

	return iterationMediaList, err
}

// CreateIterationMedia creates an iteration media in the database.
func (c *Client) CreateIterationMedia(ctx context.Context, input *models.IterationMediaCreationInput) (*models.IterationMedia, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateIterationMedia")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateIterationMedia called")

	return c.querier.CreateIterationMedia(ctx, input)
}

// UpdateIterationMedia updates a particular iteration media. Note that UpdateIterationMedia expects the
// provided input to have a valid ID.
func (c *Client) UpdateIterationMedia(ctx context.Context, updated *models.IterationMedia) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateIterationMedia")
	defer span.End()

	tracing.AttachIterationMediaIDToSpan(span, updated.ID)
	c.logger.WithValue("iteration_media_id", updated.ID).Debug("UpdateIterationMedia called")

	return c.querier.UpdateIterationMedia(ctx, updated)
}

// ArchiveIterationMedia archives an iteration media from the database by its ID.
func (c *Client) ArchiveIterationMedia(ctx context.Context, recipeIterationID, iterationMediaID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveIterationMedia")
	defer span.End()

	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
	tracing.AttachIterationMediaIDToSpan(span, iterationMediaID)

	c.logger.WithValues(map[string]interface{}{
		"iteration_media_id":  iterationMediaID,
		"recipe_iteration_id": recipeIterationID,
	}).Debug("ArchiveIterationMedia called")

	return c.querier.ArchiveIterationMedia(ctx, recipeIterationID, iterationMediaID)
}

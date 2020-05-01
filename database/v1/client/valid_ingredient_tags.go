package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.ValidIngredientTagDataManager = (*Client)(nil)

// ValidIngredientTagExists fetches whether or not a valid ingredient tag exists from the database.
func (c *Client) ValidIngredientTagExists(ctx context.Context, validIngredientTagID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "ValidIngredientTagExists")
	defer span.End()

	tracing.AttachValidIngredientTagIDToSpan(span, validIngredientTagID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_tag_id": validIngredientTagID,
	}).Debug("ValidIngredientTagExists called")

	return c.querier.ValidIngredientTagExists(ctx, validIngredientTagID)
}

// GetValidIngredientTag fetches a valid ingredient tag from the database.
func (c *Client) GetValidIngredientTag(ctx context.Context, validIngredientTagID uint64) (*models.ValidIngredientTag, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientTag")
	defer span.End()

	tracing.AttachValidIngredientTagIDToSpan(span, validIngredientTagID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_tag_id": validIngredientTagID,
	}).Debug("GetValidIngredientTag called")

	return c.querier.GetValidIngredientTag(ctx, validIngredientTagID)
}

// GetAllValidIngredientTagsCount fetches the count of valid ingredient tags from the database that meet a particular filter.
func (c *Client) GetAllValidIngredientTagsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllValidIngredientTagsCount")
	defer span.End()

	c.logger.Debug("GetAllValidIngredientTagsCount called")

	return c.querier.GetAllValidIngredientTagsCount(ctx)
}

// GetValidIngredientTags fetches a list of valid ingredient tags from the database that meet a particular filter.
func (c *Client) GetValidIngredientTags(ctx context.Context, filter *models.QueryFilter) (*models.ValidIngredientTagList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientTags")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetValidIngredientTags called")

	validIngredientTagList, err := c.querier.GetValidIngredientTags(ctx, filter)

	return validIngredientTagList, err
}

// CreateValidIngredientTag creates a valid ingredient tag in the database.
func (c *Client) CreateValidIngredientTag(ctx context.Context, input *models.ValidIngredientTagCreationInput) (*models.ValidIngredientTag, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidIngredientTag")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateValidIngredientTag called")

	return c.querier.CreateValidIngredientTag(ctx, input)
}

// UpdateValidIngredientTag updates a particular valid ingredient tag. Note that UpdateValidIngredientTag expects the
// provided input to have a valid ID.
func (c *Client) UpdateValidIngredientTag(ctx context.Context, updated *models.ValidIngredientTag) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidIngredientTag")
	defer span.End()

	tracing.AttachValidIngredientTagIDToSpan(span, updated.ID)
	c.logger.WithValue("valid_ingredient_tag_id", updated.ID).Debug("UpdateValidIngredientTag called")

	return c.querier.UpdateValidIngredientTag(ctx, updated)
}

// ArchiveValidIngredientTag archives a valid ingredient tag from the database by its ID.
func (c *Client) ArchiveValidIngredientTag(ctx context.Context, validIngredientTagID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidIngredientTag")
	defer span.End()

	tracing.AttachValidIngredientTagIDToSpan(span, validIngredientTagID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_tag_id": validIngredientTagID,
	}).Debug("ArchiveValidIngredientTag called")

	return c.querier.ArchiveValidIngredientTag(ctx, validIngredientTagID)
}

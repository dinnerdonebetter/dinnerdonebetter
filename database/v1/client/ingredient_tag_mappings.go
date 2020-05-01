package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.IngredientTagMappingDataManager = (*Client)(nil)

// IngredientTagMappingExists fetches whether or not an ingredient tag mapping exists from the database.
func (c *Client) IngredientTagMappingExists(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "IngredientTagMappingExists")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	tracing.AttachIngredientTagMappingIDToSpan(span, ingredientTagMappingID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id":       validIngredientID,
		"ingredient_tag_mapping_id": ingredientTagMappingID,
	}).Debug("IngredientTagMappingExists called")

	return c.querier.IngredientTagMappingExists(ctx, validIngredientID, ingredientTagMappingID)
}

// GetIngredientTagMapping fetches an ingredient tag mapping from the database.
func (c *Client) GetIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (*models.IngredientTagMapping, error) {
	ctx, span := tracing.StartSpan(ctx, "GetIngredientTagMapping")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	tracing.AttachIngredientTagMappingIDToSpan(span, ingredientTagMappingID)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id":       validIngredientID,
		"ingredient_tag_mapping_id": ingredientTagMappingID,
	}).Debug("GetIngredientTagMapping called")

	return c.querier.GetIngredientTagMapping(ctx, validIngredientID, ingredientTagMappingID)
}

// GetAllIngredientTagMappingsCount fetches the count of ingredient tag mappings from the database that meet a particular filter.
func (c *Client) GetAllIngredientTagMappingsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllIngredientTagMappingsCount")
	defer span.End()

	c.logger.Debug("GetAllIngredientTagMappingsCount called")

	return c.querier.GetAllIngredientTagMappingsCount(ctx)
}

// GetIngredientTagMappings fetches a list of ingredient tag mappings from the database that meet a particular filter.
func (c *Client) GetIngredientTagMappings(ctx context.Context, validIngredientID uint64, filter *models.QueryFilter) (*models.IngredientTagMappingList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetIngredientTagMappings")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"valid_ingredient_id": validIngredientID,
	}).Debug("GetIngredientTagMappings called")

	ingredientTagMappingList, err := c.querier.GetIngredientTagMappings(ctx, validIngredientID, filter)

	return ingredientTagMappingList, err
}

// CreateIngredientTagMapping creates an ingredient tag mapping in the database.
func (c *Client) CreateIngredientTagMapping(ctx context.Context, input *models.IngredientTagMappingCreationInput) (*models.IngredientTagMapping, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateIngredientTagMapping")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateIngredientTagMapping called")

	return c.querier.CreateIngredientTagMapping(ctx, input)
}

// UpdateIngredientTagMapping updates a particular ingredient tag mapping. Note that UpdateIngredientTagMapping expects the
// provided input to have a valid ID.
func (c *Client) UpdateIngredientTagMapping(ctx context.Context, updated *models.IngredientTagMapping) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateIngredientTagMapping")
	defer span.End()

	tracing.AttachIngredientTagMappingIDToSpan(span, updated.ID)
	c.logger.WithValue("ingredient_tag_mapping_id", updated.ID).Debug("UpdateIngredientTagMapping called")

	return c.querier.UpdateIngredientTagMapping(ctx, updated)
}

// ArchiveIngredientTagMapping archives an ingredient tag mapping from the database by its ID.
func (c *Client) ArchiveIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveIngredientTagMapping")
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	tracing.AttachIngredientTagMappingIDToSpan(span, ingredientTagMappingID)

	c.logger.WithValues(map[string]interface{}{
		"ingredient_tag_mapping_id": ingredientTagMappingID,
		"valid_ingredient_id":       validIngredientID,
	}).Debug("ArchiveIngredientTagMapping called")

	return c.querier.ArchiveIngredientTagMapping(ctx, validIngredientID, ingredientTagMappingID)
}

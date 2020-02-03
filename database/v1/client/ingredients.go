package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.IngredientDataManager = (*Client)(nil)

// attachIngredientIDToSpan provides a consistent way to attach an ingredient's ID to a span
func attachIngredientIDToSpan(span *trace.Span, ingredientID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("ingredient_id", strconv.FormatUint(ingredientID, 10)))
	}
}

// GetIngredient fetches an ingredient from the database
func (c *Client) GetIngredient(ctx context.Context, ingredientID, userID uint64) (*models.Ingredient, error) {
	ctx, span := trace.StartSpan(ctx, "GetIngredient")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachIngredientIDToSpan(span, ingredientID)

	c.logger.WithValues(map[string]interface{}{
		"ingredient_id": ingredientID,
		"user_id":       userID,
	}).Debug("GetIngredient called")

	return c.querier.GetIngredient(ctx, ingredientID, userID)
}

// GetIngredientCount fetches the count of ingredients from the database that meet a particular filter
func (c *Client) GetIngredientCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetIngredientCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetIngredientCount called")

	return c.querier.GetIngredientCount(ctx, filter, userID)
}

// GetAllIngredientsCount fetches the count of ingredients from the database that meet a particular filter
func (c *Client) GetAllIngredientsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllIngredientsCount")
	defer span.End()

	c.logger.Debug("GetAllIngredientsCount called")

	return c.querier.GetAllIngredientsCount(ctx)
}

// GetIngredients fetches a list of ingredients from the database that meet a particular filter
func (c *Client) GetIngredients(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.IngredientList, error) {
	ctx, span := trace.StartSpan(ctx, "GetIngredients")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetIngredients called")

	ingredientList, err := c.querier.GetIngredients(ctx, filter, userID)

	return ingredientList, err
}

// GetAllIngredientsForUser fetches a list of ingredients from the database that meet a particular filter
func (c *Client) GetAllIngredientsForUser(ctx context.Context, userID uint64) ([]models.Ingredient, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllIngredientsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllIngredientsForUser called")

	ingredientList, err := c.querier.GetAllIngredientsForUser(ctx, userID)

	return ingredientList, err
}

// CreateIngredient creates an ingredient in the database
func (c *Client) CreateIngredient(ctx context.Context, input *models.IngredientCreationInput) (*models.Ingredient, error) {
	ctx, span := trace.StartSpan(ctx, "CreateIngredient")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateIngredient called")

	return c.querier.CreateIngredient(ctx, input)
}

// UpdateIngredient updates a particular ingredient. Note that UpdateIngredient expects the
// provided input to have a valid ID.
func (c *Client) UpdateIngredient(ctx context.Context, input *models.Ingredient) error {
	ctx, span := trace.StartSpan(ctx, "UpdateIngredient")
	defer span.End()

	attachIngredientIDToSpan(span, input.ID)
	c.logger.WithValue("ingredient_id", input.ID).Debug("UpdateIngredient called")

	return c.querier.UpdateIngredient(ctx, input)
}

// ArchiveIngredient archives an ingredient from the database by its ID
func (c *Client) ArchiveIngredient(ctx context.Context, ingredientID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveIngredient")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachIngredientIDToSpan(span, ingredientID)

	c.logger.WithValues(map[string]interface{}{
		"ingredient_id": ingredientID,
		"user_id":       userID,
	}).Debug("ArchiveIngredient called")

	return c.querier.ArchiveIngredient(ctx, ingredientID, userID)
}

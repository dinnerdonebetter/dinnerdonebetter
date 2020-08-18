package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.RecipeStepProductDataManager = (*Client)(nil)

// RecipeStepProductExists fetches whether or not a recipe step product exists from the database.
func (c *Client) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepProductExists")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":              recipeID,
		"recipe_step_id":         recipeStepID,
		"recipe_step_product_id": recipeStepProductID,
	}).Debug("RecipeStepProductExists called")

	return c.querier.RecipeStepProductExists(ctx, recipeID, recipeStepID, recipeStepProductID)
}

// GetRecipeStepProduct fetches a recipe step product from the database.
func (c *Client) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*models.RecipeStepProduct, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepProduct")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":              recipeID,
		"recipe_step_id":         recipeStepID,
		"recipe_step_product_id": recipeStepProductID,
	}).Debug("GetRecipeStepProduct called")

	return c.querier.GetRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID)
}

// GetAllRecipeStepProductsCount fetches the count of recipe step products from the database that meet a particular filter.
func (c *Client) GetAllRecipeStepProductsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepProductsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepProductsCount called")

	return c.querier.GetAllRecipeStepProductsCount(ctx)
}

// GetAllRecipeStepProducts fetches a list of all recipe step products in the database.
func (c *Client) GetAllRecipeStepProducts(ctx context.Context, results chan []models.RecipeStepProduct) error {
	ctx, span := tracing.StartSpan(ctx, "GetAllRecipeStepProducts")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepProducts called")

	return c.querier.GetAllRecipeStepProducts(ctx, results)
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter.
func (c *Client) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepProductList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepProducts")
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValues(map[string]interface{}{
		"recipe_id":      recipeID,
		"recipe_step_id": recipeStepID,
	}).Debug("GetRecipeStepProducts called")

	recipeStepProductList, err := c.querier.GetRecipeStepProducts(ctx, recipeID, recipeStepID, filter)

	return recipeStepProductList, err
}

// GetRecipeStepProductsWithIDs fetches recipe step products from the database within a given set of IDs.
func (c *Client) GetRecipeStepProductsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepProduct, error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepProductsWithIDs")
	defer span.End()

	c.logger.WithValues(map[string]interface{}{
		"id_count": len(ids),
	}).Debug("GetRecipeStepProductsWithIDs called")

	recipeStepProductList, err := c.querier.GetRecipeStepProductsWithIDs(ctx, recipeID, recipeStepID, limit, ids)

	return recipeStepProductList, err
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (c *Client) CreateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProductCreationInput) (*models.RecipeStepProduct, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepProduct")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepProduct called")

	return c.querier.CreateRecipeStepProduct(ctx, input)
}

// UpdateRecipeStepProduct updates a particular recipe step product. Note that UpdateRecipeStepProduct expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepProduct(ctx context.Context, updated *models.RecipeStepProduct) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepProduct")
	defer span.End()

	tracing.AttachRecipeStepProductIDToSpan(span, updated.ID)
	c.logger.WithValue("recipe_step_product_id", updated.ID).Debug("UpdateRecipeStepProduct called")

	return c.querier.UpdateRecipeStepProduct(ctx, updated)
}

// ArchiveRecipeStepProduct archives a recipe step product from the database by its ID.
func (c *Client) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepProduct")
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_product_id": recipeStepProductID,
		"recipe_step_id":         recipeStepID,
	}).Debug("ArchiveRecipeStepProduct called")

	return c.querier.ArchiveRecipeStepProduct(ctx, recipeStepID, recipeStepProductID)
}

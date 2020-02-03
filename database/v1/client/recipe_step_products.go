package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.RecipeStepProductDataManager = (*Client)(nil)

// attachRecipeStepProductIDToSpan provides a consistent way to attach a recipe step product's ID to a span
func attachRecipeStepProductIDToSpan(span *trace.Span, recipeStepProductID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_step_product_id", strconv.FormatUint(recipeStepProductID, 10)))
	}
}

// GetRecipeStepProduct fetches a recipe step product from the database
func (c *Client) GetRecipeStepProduct(ctx context.Context, recipeStepProductID, userID uint64) (*models.RecipeStepProduct, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepProduct")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepProductIDToSpan(span, recipeStepProductID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_product_id": recipeStepProductID,
		"user_id":                userID,
	}).Debug("GetRecipeStepProduct called")

	return c.querier.GetRecipeStepProduct(ctx, recipeStepProductID, userID)
}

// GetRecipeStepProductCount fetches the count of recipe step products from the database that meet a particular filter
func (c *Client) GetRecipeStepProductCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepProductCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepProductCount called")

	return c.querier.GetRecipeStepProductCount(ctx, filter, userID)
}

// GetAllRecipeStepProductsCount fetches the count of recipe step products from the database that meet a particular filter
func (c *Client) GetAllRecipeStepProductsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepProductsCount")
	defer span.End()

	c.logger.Debug("GetAllRecipeStepProductsCount called")

	return c.querier.GetAllRecipeStepProductsCount(ctx)
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter
func (c *Client) GetRecipeStepProducts(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepProductList, error) {
	ctx, span := trace.StartSpan(ctx, "GetRecipeStepProducts")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetRecipeStepProducts called")

	recipeStepProductList, err := c.querier.GetRecipeStepProducts(ctx, filter, userID)

	return recipeStepProductList, err
}

// GetAllRecipeStepProductsForUser fetches a list of recipe step products from the database that meet a particular filter
func (c *Client) GetAllRecipeStepProductsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepProduct, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllRecipeStepProductsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllRecipeStepProductsForUser called")

	recipeStepProductList, err := c.querier.GetAllRecipeStepProductsForUser(ctx, userID)

	return recipeStepProductList, err
}

// CreateRecipeStepProduct creates a recipe step product in the database
func (c *Client) CreateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProductCreationInput) (*models.RecipeStepProduct, error) {
	ctx, span := trace.StartSpan(ctx, "CreateRecipeStepProduct")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateRecipeStepProduct called")

	return c.querier.CreateRecipeStepProduct(ctx, input)
}

// UpdateRecipeStepProduct updates a particular recipe step product. Note that UpdateRecipeStepProduct expects the
// provided input to have a valid ID.
func (c *Client) UpdateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProduct) error {
	ctx, span := trace.StartSpan(ctx, "UpdateRecipeStepProduct")
	defer span.End()

	attachRecipeStepProductIDToSpan(span, input.ID)
	c.logger.WithValue("recipe_step_product_id", input.ID).Debug("UpdateRecipeStepProduct called")

	return c.querier.UpdateRecipeStepProduct(ctx, input)
}

// ArchiveRecipeStepProduct archives a recipe step product from the database by its ID
func (c *Client) ArchiveRecipeStepProduct(ctx context.Context, recipeStepProductID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveRecipeStepProduct")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachRecipeStepProductIDToSpan(span, recipeStepProductID)

	c.logger.WithValues(map[string]interface{}{
		"recipe_step_product_id": recipeStepProductID,
		"user_id":                userID,
	}).Debug("ArchiveRecipeStepProduct called")

	return c.querier.ArchiveRecipeStepProduct(ctx, recipeStepProductID, userID)
}

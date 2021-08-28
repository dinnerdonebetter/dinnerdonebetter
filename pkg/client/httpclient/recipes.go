package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// RecipeExists retrieves whether a recipe exists.
func (c *Client) RecipeExists(ctx context.Context, recipeID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildRecipeExistsRequest(ctx, recipeID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building recipe existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for recipe #%d", recipeID)
	}

	return exists, nil
}

// GetRecipe gets a recipe.
func (c *Client) GetRecipe(ctx context.Context, recipeID uint64) (*types.FullRecipe, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildGetRecipeRequest(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get recipe request")
	}

	var recipe *types.FullRecipe
	if err = c.fetchAndUnmarshal(ctx, req, &recipe); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving recipe")
	}

	return recipe, nil
}

// GetRecipes retrieves a list of recipes.
func (c *Client) GetRecipes(ctx context.Context, filter *types.QueryFilter) (*types.RecipeList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetRecipesRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building recipes list request")
	}

	var recipes *types.RecipeList
	if err = c.fetchAndUnmarshal(ctx, req, &recipes); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving recipes")
	}

	return recipes, nil
}

// CreateRecipe creates a recipe.
func (c *Client) CreateRecipe(ctx context.Context, input *types.RecipeCreationInput) (*types.Recipe, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create recipe request")
	}

	var recipe *types.Recipe
	if err = c.fetchAndUnmarshal(ctx, req, &recipe); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating recipe")
	}

	return recipe, nil
}

// UpdateRecipe updates a recipe.
func (c *Client) UpdateRecipe(ctx context.Context, recipe *types.Recipe) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipe == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipe.ID)
	tracing.AttachRecipeIDToSpan(span, recipe.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeRequest(ctx, recipe)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update recipe request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipe); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe #%d", recipe.ID)
	}

	return nil
}

// ArchiveRecipe archives a recipe.
func (c *Client) ArchiveRecipe(ctx context.Context, recipeID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildArchiveRecipeRequest(ctx, recipeID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive recipe request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving recipe #%d", recipeID)
	}

	return nil
}

// GetAuditLogForRecipe retrieves a list of audit log entries pertaining to a recipe.
func (c *Client) GetAuditLogForRecipe(ctx context.Context, recipeID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildGetAuditLogForRecipeRequest(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for recipe request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

package httpclient

import (
	"context"

	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetRecipe gets a recipe.
func (c *Client) GetRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildGetRecipeRequest(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get recipe request")
	}

	var recipe *types.Recipe
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
func (c *Client) CreateRecipe(ctx context.Context, input *types.RecipeCreationRequestInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return "", ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return "", observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeRequest(ctx, input)
	if err != nil {
		return "", observability.PrepareError(err, logger, span, "building create recipe request")
	}

	var pwr *types.PreWriteResponse
	if err = c.fetchAndUnmarshal(ctx, req, &pwr); err != nil {
		return "", observability.PrepareError(err, logger, span, "creating recipe")
	}

	return pwr.ID, nil
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
		return observability.PrepareError(err, logger, span, "updating recipe %s", recipe.ID)
	}

	return nil
}

// ArchiveRecipe archives a recipe.
func (c *Client) ArchiveRecipe(ctx context.Context, recipeID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildArchiveRecipeRequest(ctx, recipeID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive recipe request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving recipe %s", recipeID)
	}

	return nil
}

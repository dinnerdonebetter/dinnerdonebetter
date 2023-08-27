package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetRecipe gets a recipe.
func (c *Client) GetRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildGetRecipeRequest(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe request")
	}

	var recipe *types.Recipe
	if err = c.fetchAndUnmarshal(ctx, req, &recipe); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
	}

	return recipe, nil
}

// GetRecipes retrieves a list of recipes.
func (c *Client) GetRecipes(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Recipe], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetRecipesRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipes list request")
	}

	var recipes *types.QueryFilteredResult[types.Recipe]
	if err = c.fetchAndUnmarshal(ctx, req, &recipes); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipes")
	}

	return recipes, nil
}

// SearchForRecipes retrieves a list of recipes.
func (c *Client) SearchForRecipes(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Recipe], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(c.logger.Clone())

	tracing.AttachSearchQueryToSpan(span, query)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildSearchForRecipesRequest(ctx, query, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipes list request")
	}

	var recipes *types.QueryFilteredResult[types.Recipe]
	if err = c.fetchAndUnmarshal(ctx, req, &recipes); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipes")
	}

	return recipes, nil
}

// CreateRecipe creates a recipe.
func (c *Client) CreateRecipe(ctx context.Context, input *types.RecipeCreationRequestInput) (*types.Recipe, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe request")
	}

	var recipe *types.Recipe
	if err = c.fetchAndUnmarshal(ctx, req, &recipe); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe")
	}

	return recipe, nil
}

// UpdateRecipe updates a recipe.
func (c *Client) UpdateRecipe(ctx context.Context, recipe *types.Recipe) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipe == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipe.ID)
	tracing.AttachRecipeIDToSpan(span, recipe.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeRequest(ctx, recipe)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipe); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe %s", recipe.ID)
	}

	return nil
}

// ArchiveRecipe archives a recipe.
func (c *Client) ArchiveRecipe(ctx context.Context, recipeID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildArchiveRecipeRequest(ctx, recipeID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe %s", recipeID)
	}

	return nil
}

// GetMealPlanTasksForRecipe gets a recipe.
func (c *Client) GetMealPlanTasksForRecipe(ctx context.Context, recipeID string) ([]*types.MealPlanTaskDatabaseCreationInput, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildGetRecipeMealPlanTasksRequest(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe request")
	}

	var prepSteps []*types.MealPlanTaskDatabaseCreationInput
	if err = c.fetchAndUnmarshal(ctx, req, &prepSteps); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
	}

	return prepSteps, nil
}

// UploadRecipeMedia uploads a new avatar.
func (c *Client) UploadRecipeMedia(ctx context.Context, files map[string][]byte, recipeID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return buildInvalidIDError("recipe")
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if files == nil {
		return ErrNilInputProvided
	}

	req, err := c.requestBuilder.BuildMultipleRecipeMediaUploadRequest(ctx, files, recipeID)
	if err != nil {
		return observability.PrepareError(err, span, "building media upload request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "uploading media")
	}

	return nil
}

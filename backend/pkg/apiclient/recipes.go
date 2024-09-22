package apiclient

import (
	"context"
	"image"
	"image/png"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	recipesBasePath = "recipes"
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	res, err := c.authedGeneratedClient.GetRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Recipe]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRecipes retrieves a list of recipes.
func (c *Client) GetRecipes(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Recipe], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetRecipesParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetRecipes(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "recipes list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.Recipe]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipes")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.Recipe]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// SearchForRecipes retrieves a list of recipes.
func (c *Client) SearchForRecipes(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Recipe], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(c.logger.Clone())

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	params := &generated.SearchForRecipesParams{}
	c.copyType(params, filter)
	params.Q = query

	res, err := c.authedGeneratedClient.SearchForRecipes(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "recipes list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.Recipe]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipes")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.Recipe]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
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

	body := generated.CreateRecipeJSONRequestBody{}
	c.copyType(&body, input)

	// manual body shaping
	for i, step := range input.Steps {
		for j, cc := range step.CompletionConditions {
			bodySteps := *body.Steps
			bodyCCs := bodySteps[i].CompletionConditions
			(*bodyCCs)[j].Ingredients = pointer.To(make([]int, len(input.Steps[i].CompletionConditions)))
			for k, ingredientID := range cc.Ingredients {
				(*(*bodyCCs)[j].Ingredients)[k] = int(ingredientID)
			}
		}
	}

	logger.WithValue("body", body).Info("creating recipe")

	res, err := c.authedGeneratedClient.CreateRecipe(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create recipe")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Recipe]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipe.ID)

	input := generated.UpdateRecipeJSONRequestBody{}
	c.copyType(&input, recipe)

	res, err := c.authedGeneratedClient.UpdateRecipe(ctx, recipe.ID, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update recipe")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Recipe]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	res, err := c.authedGeneratedClient.ArchiveRecipe(ctx, recipeID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive recipe")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Recipe]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	res, err := c.authedGeneratedClient.GetRecipeMealPlanTasks(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.MealPlanTaskDatabaseCreationInput]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRecipeDAG gets a recipe.
func (c *Client) GetRecipeDAG(ctx context.Context, recipeID string) (image.Image, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	res, err := c.authedGeneratedClient.GetRecipeDAG(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe")
	}
	defer c.closeResponseBody(ctx, res)

	img, err := png.Decode(res.Body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
	}

	c.closeResponseBody(ctx, res)

	return img, nil
}

// CloneRecipe gets a recipe.
func (c *Client) CloneRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	res, err := c.authedGeneratedClient.CloneRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Recipe]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UploadRecipeMedia uploads a new avatar.
// TODO: write unit test for this.
func (c *Client) UploadRecipeMedia(ctx context.Context, files map[string][]byte, recipeID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return buildInvalidIDError("recipe")
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if files == nil {
		return ErrNilInputProvided
	}

	uri := c.BuildURL(ctx, nil, recipesBasePath, recipeID, "images")

	req, err := c.buildMultipleRecipeMediaUploadRequest(ctx, uri, files)
	if err != nil {
		return observability.PrepareError(err, span, "media upload")
	}

	var apiResponse *types.APIResponse[[]*types.RecipeMedia]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "uploading media")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetRecipeStepCompletionCondition gets a recipe step completion condition.
func (c *Client) GetRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepCompletionCondition, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepIngredientID == "" {
		return nil, buildInvalidIDError("recipe step completion condition")
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepIngredientID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepIngredientID)

	res, err := c.authedGeneratedClient.GetRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe step completion condition")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.RecipeStepCompletionCondition]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step completion condition")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRecipeStepCompletionConditions retrieves a list of recipe step completion conditions.
func (c *Client) GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepCompletionCondition], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	params := &generated.GetRecipeStepCompletionConditionsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetRecipeStepCompletionConditions(ctx, recipeID, recipeStepID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "recipe step completion conditions list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.RecipeStepCompletionCondition]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step completion conditions")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.RecipeStepCompletionCondition]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateRecipeStepCompletionCondition creates a recipe step completion condition.
func (c *Client) CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput) (*types.RecipeStepCompletionCondition, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipeStep")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateRecipeStepCompletionConditionJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create recipe step completion condition")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.RecipeStepCompletionCondition]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step completion condition")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateRecipeStepCompletionCondition updates a recipe step completion condition.
func (c *Client) UpdateRecipeStepCompletionCondition(ctx context.Context, recipeID string, recipeStepIngredient *types.RecipeStepCompletionCondition) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepIngredient == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepIngredient.ID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepIngredient.ID)

	body := generated.UpdateRecipeStepCompletionConditionJSONRequestBody{}
	c.copyType(&body, recipeStepIngredient)

	res, err := c.authedGeneratedClient.UpdateRecipeStepCompletionCondition(ctx, recipeID, recipeStepIngredient.BelongsToRecipeStep, recipeStepIngredient.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update recipe step completion condition")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.RecipeStepCompletionCondition]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step completion condition")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveRecipeStepCompletionCondition archives a recipe step completion condition.
func (c *Client) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepIngredientID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepIngredientID)

	res, err := c.authedGeneratedClient.ArchiveRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive recipe step completion condition")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.RecipeStepCompletionCondition]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step completion condition")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

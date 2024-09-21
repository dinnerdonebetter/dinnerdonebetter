package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetRecipePrepTask gets a recipe prep task.
func (c *Client) GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipePrepTaskID == "" {
		return nil, buildInvalidIDError("recipe prep task")
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

	res, err := c.authedGeneratedClient.GetRecipePrepTask(ctx, recipeID, recipePrepTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe prep task")
	}

	var apiResponse *types.APIResponse[*types.RecipePrepTask]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe prep task")
	}

	return apiResponse.Data, nil
}

// GetRecipePrepTasks retrieves a list of recipe prep tasks.
func (c *Client) GetRecipePrepTasks(ctx context.Context, recipeID string, filter *types.QueryFilter) ([]*types.RecipePrepTask, error) {
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

	params := &generated.GetRecipePrepTasksParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetRecipePrepTasks(ctx, recipeID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "recipe prep tasks list")
	}

	var apiResponse *types.APIResponse[[]*types.RecipePrepTask]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe prep tasks")
	}

	return apiResponse.Data, nil
}

// CreateRecipePrepTask creates a recipe prep task.
func (c *Client) CreateRecipePrepTask(ctx context.Context, input *types.RecipePrepTaskCreationRequestInput) (*types.RecipePrepTask, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateRecipePrepTaskJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateRecipePrepTask(ctx, input.BelongsToRecipe, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create recipe prep task")
	}

	var apiResponse *types.APIResponse[*types.RecipePrepTask]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
	}

	return apiResponse.Data, nil
}

// UpdateRecipePrepTask updates a recipe prep task.
func (c *Client) UpdateRecipePrepTask(ctx context.Context, recipePrepTask *types.RecipePrepTask) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipePrepTask == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTask.ID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTask.ID)

	body := generated.UpdateRecipePrepTaskJSONRequestBody{}
	c.copyType(&body, recipePrepTask)

	res, err := c.authedGeneratedClient.UpdateRecipePrepTask(ctx, recipePrepTask.BelongsToRecipe, recipePrepTask.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update recipe prep task")
	}

	var apiResponse *types.APIResponse[*types.RecipePrepTask]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task")
	}

	return nil
}

// ArchiveRecipePrepTask archives a recipe prep task.
func (c *Client) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipePrepTaskID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

	res, err := c.authedGeneratedClient.ArchiveRecipePrepTask(ctx, recipeID, recipePrepTaskID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive recipe prep task")
	}

	var apiResponse *types.APIResponse[*types.RecipePrepTask]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe prep task")
	}

	return nil
}

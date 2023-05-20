package apiclient

import (
	"context"

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
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipePrepTaskID == "" {
		return nil, buildInvalidIDError("recipe prep task")
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	req, err := c.requestBuilder.BuildGetRecipePrepTaskRequest(ctx, recipeID, recipePrepTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe prep task request")
	}

	var recipePrepTask *types.RecipePrepTask
	if err = c.fetchAndUnmarshal(ctx, req, &recipePrepTask); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe prep task")
	}

	return recipePrepTask, nil
}

// GetRecipePrepTasks retrieves a list of recipe prep tasks.
func (c *Client) GetRecipePrepTasks(ctx context.Context, recipeID string, filter *types.QueryFilter) ([]*types.RecipePrepTask, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildGetRecipePrepTasksRequest(ctx, recipeID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe prep tasks list request")
	}

	var recipePrepTasks []*types.RecipePrepTask
	if err = c.fetchAndUnmarshal(ctx, req, &recipePrepTasks); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe prep tasks")
	}

	return recipePrepTasks, nil
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

	req, err := c.requestBuilder.BuildCreateRecipePrepTaskRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe prep task request")
	}

	var recipePrepTask *types.RecipePrepTask
	if err = c.fetchAndUnmarshal(ctx, req, &recipePrepTask); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
	}

	return recipePrepTask, nil
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
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTask.ID)

	req, err := c.requestBuilder.BuildUpdateRecipePrepTaskRequest(ctx, recipePrepTask)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe prep task request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipePrepTask); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task %s", recipePrepTask.ID)
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
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipePrepTaskID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	req, err := c.requestBuilder.BuildArchiveRecipePrepTaskRequest(ctx, recipeID, recipePrepTaskID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe prep task request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe prep task %s", recipePrepTaskID)
	}

	return nil
}

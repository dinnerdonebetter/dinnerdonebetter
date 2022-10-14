package apiclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetRecipePrepTask gets a recipe step.
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
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	req, err := c.requestBuilder.BuildGetRecipePrepTaskRequest(ctx, recipeID, recipePrepTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe step request")
	}

	var recipePrepTask *types.RecipePrepTask
	if err = c.fetchAndUnmarshal(ctx, req, &recipePrepTask); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step")
	}

	return recipePrepTask, nil
}

// GetRecipePrepTasks retrieves a list of recipe steps.
func (c *Client) GetRecipePrepTasks(ctx context.Context, recipeID string, filter *types.QueryFilter) (*types.RecipePrepTaskList, error) {
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
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe steps list request")
	}

	var recipePrepTasks *types.RecipePrepTaskList
	if err = c.fetchAndUnmarshal(ctx, req, &recipePrepTasks); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe steps")
	}

	return recipePrepTasks, nil
}

// CreateRecipePrepTask creates a recipe step.
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
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe step request")
	}

	var recipePrepTask *types.RecipePrepTask
	if err = c.fetchAndUnmarshal(ctx, req, &recipePrepTask); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step")
	}

	return recipePrepTask, nil
}

// UpdateRecipePrepTask updates a recipe step.
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
		return observability.PrepareAndLogError(err, logger, span, "building update recipe step request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipePrepTask); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step %s", recipePrepTask.ID)
	}

	return nil
}

// ArchiveRecipePrepTask archives a recipe step.
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
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe step request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step %s", recipePrepTaskID)
	}

	return nil
}

package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

// GetRecipeStep gets a recipe step.
func (c *Client) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	req, err := c.requestBuilder.BuildGetRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe step request")
	}

	var recipeStep *types.RecipeStep
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStep); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step")
	}

	return recipeStep, nil
}

// GetRecipeSteps retrieves a list of recipe steps.
func (c *Client) GetRecipeSteps(ctx context.Context, recipeID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStep], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildGetRecipeStepsRequest(ctx, recipeID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe steps list request")
	}

	var recipeSteps *types.QueryFilteredResult[types.RecipeStep]
	if err = c.fetchAndUnmarshal(ctx, req, &recipeSteps); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe steps")
	}

	return recipeSteps, nil
}

// CreateRecipeStep creates a recipe step.
func (c *Client) CreateRecipeStep(ctx context.Context, input *types.RecipeStepCreationRequestInput) (*types.RecipeStep, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe step request")
	}

	var recipeStep *types.RecipeStep
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStep); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step")
	}

	return recipeStep, nil
}

// UpdateRecipeStep updates a recipe step.
func (c *Client) UpdateRecipeStep(ctx context.Context, recipeStep *types.RecipeStep) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeStep == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStep.ID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStep.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepRequest(ctx, recipeStep)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe step request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStep); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step %s", recipeStep.ID)
	}

	return nil
}

// ArchiveRecipeStep archives a recipe step.
func (c *Client) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive recipe step request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "archiving recipe step %s", recipeStepID)
	}

	return nil
}

// UploadRecipeMediaForStep uploads a new avatar.
func (c *Client) UploadRecipeMediaForStep(ctx context.Context, files map[string][]byte, recipeID, recipeStepID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if files == nil {
		return ErrNilInputProvided
	}

	req, err := c.requestBuilder.BuildMultipleRecipeMediaUploadRequestForRecipeStep(ctx, files, recipeID, recipeStepID)
	if err != nil {
		return observability.PrepareError(err, span, "building media upload request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "uploading media")
	}

	return nil
}

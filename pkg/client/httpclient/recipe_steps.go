package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// RecipeStepExists retrieves whether a recipe step exists.
func (c *Client) RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	req, err := c.requestBuilder.BuildRecipeStepExistsRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building recipe step existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for recipe step #%d", recipeStepID)
	}

	return exists, nil
}

// GetRecipeStep gets a recipe step.
func (c *Client) GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*types.RecipeStep, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	req, err := c.requestBuilder.BuildGetRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get recipe step request")
	}

	var recipeStep *types.RecipeStep
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStep); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving recipe step")
	}

	return recipeStep, nil
}

// GetRecipeSteps retrieves a list of recipe steps.
func (c *Client) GetRecipeSteps(ctx context.Context, recipeID uint64, filter *types.QueryFilter) (*types.RecipeStepList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	req, err := c.requestBuilder.BuildGetRecipeStepsRequest(ctx, recipeID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building recipe steps list request")
	}

	var recipeSteps *types.RecipeStepList
	if err = c.fetchAndUnmarshal(ctx, req, &recipeSteps); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving recipe steps")
	}

	return recipeSteps, nil
}

// CreateRecipeStep creates a recipe step.
func (c *Client) CreateRecipeStep(ctx context.Context, input *types.RecipeStepCreationInput) (*types.RecipeStep, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create recipe step request")
	}

	var recipeStep *types.RecipeStep
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStep); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating recipe step")
	}

	return recipeStep, nil
}

// UpdateRecipeStep updates a recipe step.
func (c *Client) UpdateRecipeStep(ctx context.Context, recipeStep *types.RecipeStep) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeStep == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStep.ID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStep.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepRequest(ctx, recipeStep)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update recipe step request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStep); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step #%d", recipeStep.ID)
	}

	return nil
}

// ArchiveRecipeStep archives a recipe step.
func (c *Client) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive recipe step request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving recipe step #%d", recipeStepID)
	}

	return nil
}

// GetAuditLogForRecipeStep retrieves a list of audit log entries pertaining to a recipe step.
func (c *Client) GetAuditLogForRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	req, err := c.requestBuilder.BuildGetAuditLogForRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for recipe step request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

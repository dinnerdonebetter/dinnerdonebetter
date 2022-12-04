package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
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
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepIngredientID == "" {
		return nil, buildInvalidIDError("recipe step completion condition")
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildGetRecipeStepCompletionConditionRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe step completion condition request")
	}

	var recipeStepIngredient *types.RecipeStepCompletionCondition
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step completion condition")
	}

	return recipeStepIngredient, nil
}

// GetRecipeStepCompletionConditions retrieves a list of recipe step completion conditions.
func (c *Client) GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepCompletionCondition], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

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

	req, err := c.requestBuilder.BuildGetRecipeStepCompletionConditionsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe step completion conditions list request")
	}

	var recipeStepIngredients *types.QueryFilteredResult[types.RecipeStepCompletionCondition]
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredients); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step completion conditions")
	}

	return recipeStepIngredients, nil
}

// CreateRecipeStepCompletionCondition creates a recipe step completion condition.
func (c *Client) CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionCreationRequestInput) (*types.RecipeStepCompletionCondition, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipeStep")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepCompletionConditionRequest(ctx, recipeID, recipeStepID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe step completion condition request")
	}

	var recipeStepIngredient *types.RecipeStepCompletionCondition
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step completion condition")
	}

	return recipeStepIngredient, nil
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
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepIngredient == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepIngredient.ID)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepIngredient.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepCompletionConditionRequest(ctx, recipeID, recipeStepIngredient)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe step completion condition request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step completion condition %s", recipeStepIngredient.ID)
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
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepCompletionConditionRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe step completion condition request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step completion condition %s", recipeStepIngredientID)
	}

	return nil
}

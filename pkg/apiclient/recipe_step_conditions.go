package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

// GetRecipeStepCondition gets a recipe step condition.
func (c *Client) GetRecipeStepCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepCondition, error) {
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
		return nil, buildInvalidIDError("recipe step condition")
	}
	logger = logger.WithValue(keys.RecipeStepConditionIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildGetRecipeStepConditionRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe step condition request")
	}

	var recipeStepIngredient *types.RecipeStepCondition
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step condition")
	}

	return recipeStepIngredient, nil
}

// GetRecipeStepConditions retrieves a list of recipe step conditions.
func (c *Client) GetRecipeStepConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepCondition], error) {
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

	req, err := c.requestBuilder.BuildGetRecipeStepConditionsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe step conditions list request")
	}

	var recipeStepIngredients *types.QueryFilteredResult[types.RecipeStepCondition]
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredients); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step conditions")
	}

	return recipeStepIngredients, nil
}

// CreateRecipeStepCondition creates a recipe step condition.
func (c *Client) CreateRecipeStepCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepConditionCreationRequestInput) (*types.RecipeStepCondition, error) {
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

	req, err := c.requestBuilder.BuildCreateRecipeStepConditionRequest(ctx, recipeID, recipeStepID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe step condition request")
	}

	var recipeStepIngredient *types.RecipeStepCondition
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step condition")
	}

	return recipeStepIngredient, nil
}

// UpdateRecipeStepCondition updates a recipe step condition.
func (c *Client) UpdateRecipeStepCondition(ctx context.Context, recipeID string, recipeStepIngredient *types.RecipeStepCondition) error {
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
	logger = logger.WithValue(keys.RecipeStepConditionIDKey, recipeStepIngredient.ID)
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepIngredient.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepConditionRequest(ctx, recipeID, recipeStepIngredient)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe step condition request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step condition %s", recipeStepIngredient.ID)
	}

	return nil
}

// ArchiveRecipeStepCondition archives a recipe step condition.
func (c *Client) ArchiveRecipeStepCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error {
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
	logger = logger.WithValue(keys.RecipeStepConditionIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepConditionRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe step condition request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step condition %s", recipeStepIngredientID)
	}

	return nil
}

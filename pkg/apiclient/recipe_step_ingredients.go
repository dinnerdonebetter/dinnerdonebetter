package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

// GetRecipeStepIngredient gets a recipe step ingredient.
func (c *Client) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error) {
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
		return nil, buildInvalidIDError("recipe step ingredient")
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildGetRecipeStepIngredientRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe step ingredient request")
	}

	var recipeStepIngredient *types.RecipeStepIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step ingredient")
	}

	return recipeStepIngredient, nil
}

// GetRecipeStepIngredients retrieves a list of recipe step ingredients.
func (c *Client) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.RecipeStepIngredientList, error) {
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

	req, err := c.requestBuilder.BuildGetRecipeStepIngredientsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe step ingredients list request")
	}

	var recipeStepIngredients *types.RecipeStepIngredientList
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredients); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step ingredients")
	}

	return recipeStepIngredients, nil
}

// CreateRecipeStepIngredient creates a recipe step ingredient.
func (c *Client) CreateRecipeStepIngredient(ctx context.Context, recipeID string, input *types.RecipeStepIngredientCreationRequestInput) (*types.RecipeStepIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepIngredientRequest(ctx, recipeID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe step ingredient request")
	}

	var recipeStepIngredient *types.RecipeStepIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step ingredient")
	}

	return recipeStepIngredient, nil
}

// UpdateRecipeStepIngredient updates a recipe step ingredient.
func (c *Client) UpdateRecipeStepIngredient(ctx context.Context, recipeID string, recipeStepIngredient *types.RecipeStepIngredient) error {
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
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredient.ID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredient.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepIngredientRequest(ctx, recipeID, recipeStepIngredient)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe step ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step ingredient %s", recipeStepIngredient.ID)
	}

	return nil
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient.
func (c *Client) ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error {
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
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepIngredientRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe step ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step ingredient %s", recipeStepIngredientID)
	}

	return nil
}

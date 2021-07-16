package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// RecipeStepIngredientExists retrieves whether a recipe step ingredient exists.
func (c *Client) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (bool, error) {
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

	if recipeStepIngredientID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildRecipeStepIngredientExistsRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building recipe step ingredient existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for recipe step ingredient #%d", recipeStepIngredientID)
	}

	return exists, nil
}

// GetRecipeStepIngredient gets a recipe step ingredient.
func (c *Client) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*types.RecipeStepIngredient, error) {
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

	if recipeStepIngredientID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildGetRecipeStepIngredientRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get recipe step ingredient request")
	}

	var recipeStepIngredient *types.RecipeStepIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving recipe step ingredient")
	}

	return recipeStepIngredient, nil
}

// GetRecipeStepIngredients retrieves a list of recipe step ingredients.
func (c *Client) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *types.QueryFilter) (*types.RecipeStepIngredientList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

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

	req, err := c.requestBuilder.BuildGetRecipeStepIngredientsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building recipe step ingredients list request")
	}

	var recipeStepIngredients *types.RecipeStepIngredientList
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredients); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving recipe step ingredients")
	}

	return recipeStepIngredients, nil
}

// CreateRecipeStepIngredient creates a recipe step ingredient.
func (c *Client) CreateRecipeStepIngredient(ctx context.Context, recipeID uint64, input *types.RecipeStepIngredientCreationInput) (*types.RecipeStepIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepIngredientRequest(ctx, recipeID, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create recipe step ingredient request")
	}

	var recipeStepIngredient *types.RecipeStepIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating recipe step ingredient")
	}

	return recipeStepIngredient, nil
}

// UpdateRecipeStepIngredient updates a recipe step ingredient.
func (c *Client) UpdateRecipeStepIngredient(ctx context.Context, recipeID uint64, recipeStepIngredient *types.RecipeStepIngredient) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if recipeID == 0 {
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
		return observability.PrepareError(err, logger, span, "building update recipe step ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepIngredient); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step ingredient #%d", recipeStepIngredient.ID)
	}

	return nil
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient.
func (c *Client) ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) error {
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

	if recipeStepIngredientID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepIngredientRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive recipe step ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving recipe step ingredient #%d", recipeStepIngredientID)
	}

	return nil
}

// GetAuditLogForRecipeStepIngredient retrieves a list of audit log entries pertaining to a recipe step ingredient.
func (c *Client) GetAuditLogForRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) ([]*types.AuditLogEntry, error) {
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

	if recipeStepIngredientID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	req, err := c.requestBuilder.BuildGetAuditLogForRecipeStepIngredientRequest(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for recipe step ingredient request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

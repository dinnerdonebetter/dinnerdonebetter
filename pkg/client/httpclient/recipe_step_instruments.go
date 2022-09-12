package httpclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetRecipeStepInstrument gets a recipe step instrument.
func (c *Client) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
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

	if recipeStepInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	req, err := c.requestBuilder.BuildGetRecipeStepInstrumentRequest(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building get recipe step instrument request")
	}

	var recipeStepInstrument *types.RecipeStepInstrument
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepInstrument); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving recipe step instrument")
	}

	return recipeStepInstrument, nil
}

// GetRecipeStepInstruments retrieves a list of recipe step instruments.
func (c *Client) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.RecipeStepInstrumentList, error) {
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

	req, err := c.requestBuilder.BuildGetRecipeStepInstrumentsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building recipe step instruments list request")
	}

	var recipeStepInstruments *types.RecipeStepInstrumentList
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepInstruments); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving recipe step instruments")
	}

	return recipeStepInstruments, nil
}

// CreateRecipeStepInstrument creates a recipe step instrument.
func (c *Client) CreateRecipeStepInstrument(ctx context.Context, recipeID string, input *types.RecipeStepInstrumentCreationRequestInput) (*types.RecipeStepInstrument, error) {
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
		return nil, observability.PrepareError(err, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepInstrumentRequest(ctx, recipeID, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building create recipe step instrument request")
	}

	var recipeStepInstrument *types.RecipeStepInstrument
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepInstrument); err != nil {
		return nil, observability.PrepareError(err, span, "creating recipe step instrument")
	}

	return recipeStepInstrument, nil
}

// UpdateRecipeStepInstrument updates a recipe step instrument.
func (c *Client) UpdateRecipeStepInstrument(ctx context.Context, recipeID string, recipeStepInstrument *types.RecipeStepInstrument) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepInstrument == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrument.ID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrument.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepInstrumentRequest(ctx, recipeID, recipeStepInstrument)
	if err != nil {
		return observability.PrepareError(err, span, "building update recipe step instrument request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepInstrument); err != nil {
		return observability.PrepareError(err, span, "updating recipe step instrument %s", recipeStepInstrument.ID)
	}

	return nil
}

// ArchiveRecipeStepInstrument archives a recipe step instrument.
func (c *Client) ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error {
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

	if recipeStepInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepInstrumentRequest(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive recipe step instrument request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "archiving recipe step instrument %s", recipeStepInstrumentID)
	}

	return nil
}

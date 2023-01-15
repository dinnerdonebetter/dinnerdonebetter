package apiclient

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

// GetRecipeStepVessel gets a recipe step vessel.
func (c *Client) GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepVessel, error) {
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
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepInstrumentID)
	tracing.AttachRecipeStepVesselIDToSpan(span, recipeStepInstrumentID)

	req, err := c.requestBuilder.BuildGetRecipeStepVesselRequest(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe step vessel request")
	}

	var recipeStepInstrument *types.RecipeStepVessel
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepInstrument); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step vessel")
	}

	return recipeStepInstrument, nil
}

// GetRecipeStepVessels retrieves a list of recipe step vessels.
func (c *Client) GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepVessel], error) {
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

	req, err := c.requestBuilder.BuildGetRecipeStepVesselsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe step vessels list request")
	}

	var recipeStepInstruments *types.QueryFilteredResult[types.RecipeStepVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step vessels")
	}

	return recipeStepInstruments, nil
}

// CreateRecipeStepVessel creates a recipe step vessel.
func (c *Client) CreateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselCreationRequestInput) (*types.RecipeStepVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepVesselRequest(ctx, recipeID, recipeStepID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe step vessel request")
	}

	var recipeStepInstrument *types.RecipeStepVessel
	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepInstrument); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step vessel")
	}

	return recipeStepInstrument, nil
}

// UpdateRecipeStepVessel updates a recipe step vessel.
func (c *Client) UpdateRecipeStepVessel(ctx context.Context, recipeID string, recipeStepInstrument *types.RecipeStepVessel) error {
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
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepInstrument.ID)
	tracing.AttachRecipeStepVesselIDToSpan(span, recipeStepInstrument.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepVesselRequest(ctx, recipeID, recipeStepInstrument)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe step vessel request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeStepInstrument); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step vessel %s", recipeStepInstrument.ID)
	}

	return nil
}

// ArchiveRecipeStepVessel archives a recipe step vessel.
func (c *Client) ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error {
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
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepInstrumentID)
	tracing.AttachRecipeStepVesselIDToSpan(span, recipeStepInstrumentID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepVesselRequest(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe step vessel request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step vessel %s", recipeStepInstrumentID)
	}

	return nil
}

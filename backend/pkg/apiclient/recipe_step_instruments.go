package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	res, err := c.authedGeneratedClient.GetRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe step instrument")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step instrument")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRecipeStepInstruments retrieves a list of recipe step instruments.
func (c *Client) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepInstrument], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	params := &generated.GetRecipeStepInstrumentsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetRecipeStepInstruments(ctx, recipeID, recipeStepID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "recipe step instruments list")
	}

	var apiResponse *types.APIResponse[[]*types.RecipeStepInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step instruments")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.RecipeStepInstrument]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateRecipeStepInstrument creates a recipe step instrument.
func (c *Client) CreateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentCreationRequestInput) (*types.RecipeStepInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateRecipeStepInstrumentJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateRecipeStepInstrument(ctx, recipeID, recipeStepID, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create recipe step instrument")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step instrument")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepInstrument == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrument.ID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrument.ID)

	body := generated.UpdateRecipeStepInstrumentJSONRequestBody{}
	c.copyType(&body, recipeStepInstrument)

	res, err := c.authedGeneratedClient.UpdateRecipeStepInstrument(ctx, recipeID, recipeStepInstrument.BelongsToRecipeStep, recipeStepInstrument.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update recipe step instrument")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step instrument %s", recipeStepInstrument.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
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
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	res, err := c.authedGeneratedClient.ArchiveRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive recipe step instrument")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step instrument %s", recipeStepInstrumentID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

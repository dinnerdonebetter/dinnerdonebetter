package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetRecipeStepVessel gets a recipe step vessel.
func (c *Client) GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error) {
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

	if recipeStepVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVesselID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)

	res, err := c.authedGeneratedClient.GetRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe step vessel")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step vessel")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRecipeStepVessels retrieves a list of recipe step vessels.
func (c *Client) GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepVessel], error) {
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

	params := &generated.GetRecipeStepVesselsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetRecipeStepVessels(ctx, recipeID, recipeStepID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "recipe step vessels list")
	}

	var apiResponse *types.APIResponse[[]*types.RecipeStepVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step vessels")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.RecipeStepVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
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

	body := generated.CreateRecipeStepVesselJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateRecipeStepVessel(ctx, recipeID, recipeStepID, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create recipe step vessel")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step vessel")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateRecipeStepVessel updates a recipe step vessel.
func (c *Client) UpdateRecipeStepVessel(ctx context.Context, recipeID string, recipeStepVessel *types.RecipeStepVessel) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepVessel == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVessel.ID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVessel.ID)

	body := generated.UpdateRecipeStepVesselJSONRequestBody{}
	c.copyType(&body, recipeStepVessel)

	res, err := c.authedGeneratedClient.UpdateRecipeStepVessel(ctx, recipeID, recipeStepVessel.BelongsToRecipeStep, recipeStepVessel.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update recipe step vessel")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step vessel %s", recipeStepVessel.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveRecipeStepVessel archives a recipe step vessel.
func (c *Client) ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) error {
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

	if recipeStepVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVesselID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)

	res, err := c.authedGeneratedClient.ArchiveRecipeStepVessel(ctx, recipeID, recipeStepID, recipeStepVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive recipe step vessel")
	}

	var apiResponse *types.APIResponse[*types.RecipeStepVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step vessel %s", recipeStepVesselID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

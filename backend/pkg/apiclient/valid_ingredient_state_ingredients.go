package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidIngredientStateIngredient gets a valid ingredient state ingredient.
func (c *Client) GetValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	res, err := c.authedGeneratedClient.GetValidIngredientStateIngredient(ctx, validIngredientStateIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get valid ingredient state ingredient")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientStateIngredient]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidIngredientStateIngredients retrieves a list of valid ingredient state ingredient.
func (c *Client) GetValidIngredientStateIngredients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidIngredientStateIngredientsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientStateIngredients(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "valid ingredient state ingredient list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientStateIngredient]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidIngredientStateIngredientsForIngredient retrieves a list of valid ingredient state ingredient.
func (c *Client) GetValidIngredientStateIngredientsForIngredient(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	params := &generated.GetValidIngredientStateIngredientsByIngredientParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientStateIngredientsByIngredient(ctx, validIngredientID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "valid ingredient state ingredient list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientStateIngredient]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidIngredientStateIngredientsForIngredientState retrieves a list of valid ingredient state ingredient.
func (c *Client) GetValidIngredientStateIngredientsForIngredientState(ctx context.Context, ingredientState string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if ingredientState == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, ingredientState)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, ingredientState)

	params := &generated.GetValidIngredientStateIngredientsByIngredientStateParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientStateIngredientsByIngredientState(ctx, ingredientState, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "valid ingredient state ingredient list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientStateIngredient]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateValidIngredientStateIngredient creates a valid ingredient state ingredient.
func (c *Client) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientCreationRequestInput) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidIngredientStateIngredientJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidIngredientStateIngredient(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create valid ingredient state ingredient")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientStateIngredient]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidIngredientStateIngredient updates a valid ingredient state ingredient.
func (c *Client) UpdateValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredient *types.ValidIngredientStateIngredient) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateIngredient == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredient.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredient.ID)

	body := generated.UpdateValidIngredientStateIngredientJSONRequestBody{}
	c.copyType(&body, validIngredientStateIngredient)
	body.ValidIngredientID = &validIngredientStateIngredient.Ingredient.ID
	body.ValidIngredientStateID = &validIngredientStateIngredient.IngredientState.ID

	res, err := c.authedGeneratedClient.UpdateValidIngredientStateIngredient(ctx, validIngredientStateIngredient.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update valid ingredient state ingredient")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientStateIngredient]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state ingredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidIngredientStateIngredient archives a valid ingredient state ingredient.
func (c *Client) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	res, err := c.authedGeneratedClient.ArchiveValidIngredientStateIngredient(ctx, validIngredientStateIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive valid ingredient state ingredient")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientStateIngredient]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient state ingredient")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

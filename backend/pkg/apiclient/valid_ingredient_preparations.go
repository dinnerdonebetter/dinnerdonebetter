package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidIngredientPreparation gets a valid ingredient preparation.
func (c *Client) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	res, err := c.authedGeneratedClient.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidIngredientPreparations retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientPreparations(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientPreparation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidIngredientPreparationsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientPreparations(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient preparations list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientPreparation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidIngredientPreparationsForIngredient retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientPreparationsForIngredient(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientPreparation], error) {
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

	params := &generated.GetValidIngredientPreparationsByIngredientParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientPreparationsByIngredient(ctx, validIngredientID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient preparations by ingredient list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientPreparation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidIngredientPreparationsForPreparation retrieves a list of valid ingredient preparations.
func (c *Client) GetValidIngredientPreparationsForPreparation(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientPreparation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	params := &generated.GetValidIngredientPreparationsByPreparationParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientPreparationsByPreparation(ctx, validPreparationID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient preparations by preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientPreparation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidIngredientPreparationsForPreparationAndIngredientName retrieves a list of valid ingredient preparations for a given preparation ID and a given ingredient name query.
func (c *Client) GetValidIngredientPreparationsForPreparationAndIngredientName(ctx context.Context, validPreparationID, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientPreparation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	params := &generated.GetValidIngredientsByPreparationParams{}
	c.copyType(params, filter)
	params.Q = query

	res, err := c.authedGeneratedClient.GetValidIngredientsByPreparation(ctx, validPreparationID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient preparations by preparation and ingredient name")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientPreparation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateValidIngredientPreparation creates a valid ingredient preparation.
func (c *Client) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationRequestInput) (*types.ValidIngredientPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidIngredientPreparationJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidIngredientPreparation(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidIngredientPreparation updates a valid ingredient preparation.
func (c *Client) UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparation *types.ValidIngredientPreparation) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientPreparation == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparation.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparation.ID)

	body := generated.UpdateValidIngredientPreparationJSONRequestBody{}
	c.copyType(&body, validIngredientPreparation)

	res, err := c.authedGeneratedClient.UpdateValidIngredientPreparation(ctx, validIngredientPreparation.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidIngredientPreparation archives a valid ingredient preparation.
func (c *Client) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	res, err := c.authedGeneratedClient.ArchiveValidIngredientPreparation(ctx, validIngredientPreparationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

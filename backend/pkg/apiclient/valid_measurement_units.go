package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidMeasurementUnit gets a valid measurement unit.
func (c *Client) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	res, err := c.authedGeneratedClient.GetValidMeasurementUnit(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get valid measurement unit")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// SearchValidMeasurementUnits searches through a list of valid measurement units.
func (c *Client) SearchValidMeasurementUnits(ctx context.Context, query string, limit uint8) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	if limit == 0 {
		limit = types.DefaultQueryFilterLimit
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	params := &generated.SearchForValidMeasurementUnitsParams{
		Q:     query,
		Limit: int(limit),
	}

	res, err := c.authedGeneratedClient.SearchForValidMeasurementUnits(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "search for valid measurement units")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement units")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// SearchValidMeasurementUnitsByIngredientID searches through a list of valid measurement units.
func (c *Client) SearchValidMeasurementUnitsByIngredientID(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.logger.Clone()
	tracing.AttachQueryFilterToSpan(span, filter)

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	params := &generated.SearchValidMeasurementUnitsByIngredientParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.SearchValidMeasurementUnitsByIngredient(ctx, validIngredientID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "search for valid measurement units")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement units")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidMeasurementUnit]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidMeasurementUnits retrieves a list of valid measurement units.
func (c *Client) GetValidMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidMeasurementUnitsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidMeasurementUnits(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "valid measurement units list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid measurement units")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidMeasurementUnit]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateValidMeasurementUnit creates a valid measurement unit.
func (c *Client) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitCreationRequestInput) (*types.ValidMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidMeasurementUnitJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidMeasurementUnit(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create valid measurement unit")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidMeasurementUnit updates a valid measurement unit.
func (c *Client) UpdateValidMeasurementUnit(ctx context.Context, validMeasurementUnit *types.ValidMeasurementUnit) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnit == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnit.ID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnit.ID)

	body := generated.UpdateValidMeasurementUnitJSONRequestBody{}
	c.copyType(&body, validMeasurementUnit)

	res, err := c.authedGeneratedClient.UpdateValidMeasurementUnit(ctx, validMeasurementUnit.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update valid measurement unit")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement unit %s", validMeasurementUnit.ID)
	}

	return nil
}

// ArchiveValidMeasurementUnit archives a valid measurement unit.
func (c *Client) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	res, err := c.authedGeneratedClient.ArchiveValidMeasurementUnit(ctx, validMeasurementUnitID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive valid measurement unit")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement unit %s", validMeasurementUnitID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

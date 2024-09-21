package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidIngredientMeasurementUnit gets a valid ingredient measurement unit.
func (c *Client) GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	res, err := c.authedGeneratedClient.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient measurement unit request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidIngredientMeasurementUnits retrieves a list of valid ingredient measurement unit.
func (c *Client) GetValidIngredientMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientMeasurementUnit], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidIngredientMeasurementUnitsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientMeasurementUnits(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient measurement unit list request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidIngredientMeasurementUnitsForIngredient retrieves a list of valid ingredient measurement unit.
func (c *Client) GetValidIngredientMeasurementUnitsForIngredient(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientMeasurementUnit], error) {
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

	params := &generated.GetValidIngredientMeasurementUnitsByIngredientParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientMeasurementUnitsByIngredient(ctx, validIngredientID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient measurement unit list request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidIngredientMeasurementUnitsForMeasurementUnit retrieves a list of valid ingredient measurement unit.
func (c *Client) GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx context.Context, validMeasurementUnitID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientMeasurementUnit], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validMeasurementUnitID)

	params := &generated.GetValidIngredientMeasurementUnitsByMeasurementUnitParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx, validMeasurementUnitID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient measurement unit list request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateValidIngredientMeasurementUnit creates a valid ingredient measurement unit.
func (c *Client) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitCreationRequestInput) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidIngredientMeasurementUnitJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidIngredientMeasurementUnit(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient measurement unit request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidIngredientMeasurementUnit updates a valid ingredient measurement unit.
func (c *Client) UpdateValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientMeasurementUnit == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnit.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnit.ID)

	body := generated.UpdateValidIngredientMeasurementUnitJSONRequestBody{}
	c.copyType(&body, validIngredientMeasurementUnit)

	res, err := c.authedGeneratedClient.UpdateValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnit.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient measurement unit request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidIngredientMeasurementUnit archives a valid ingredient measurement unit.
func (c *Client) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	res, err := c.authedGeneratedClient.ArchiveValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient measurement unit request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientMeasurementUnit]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient measurement unit")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

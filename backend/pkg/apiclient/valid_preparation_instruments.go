package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidPreparationInstrument gets a valid ingredient preparation.
func (c *Client) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	res, err := c.authedGeneratedClient.GetValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient preparation request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparationInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient preparation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidPreparationInstruments retrieves a list of valid preparation instruments.
func (c *Client) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationInstrument], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidPreparationInstrumentsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidPreparationInstruments(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation instruments list request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidPreparationInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation instruments")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidPreparationInstrumentsForPreparation retrieves a list of valid preparation instruments.
func (c *Client) GetValidPreparationInstrumentsForPreparation(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationInstrument], error) {
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

	params := &generated.GetValidPreparationInstrumentsByPreparationParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidPreparationInstrumentsByPreparation(ctx, validPreparationID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation instruments list request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidPreparationInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation instruments")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetValidPreparationInstrumentsForInstrument retrieves a list of valid preparation instruments.
func (c *Client) GetValidPreparationInstrumentsForInstrument(ctx context.Context, validInstrumentID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationInstrument], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	params := &generated.GetValidPreparationInstrumentsByInstrumentParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidPreparationInstrumentsByInstrument(ctx, validInstrumentID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation instruments list request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidPreparationInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation instruments")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateValidPreparationInstrument creates a valid ingredient preparation.
func (c *Client) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*types.ValidPreparationInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidPreparationInstrumentJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidPreparationInstrument(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient preparation request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparationInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient preparation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidPreparationInstrument updates a valid ingredient preparation.
func (c *Client) UpdateValidPreparationInstrument(ctx context.Context, validPreparationInstrument *types.ValidPreparationInstrument) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationInstrument == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrument.ID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrument.ID)

	body := generated.UpdateValidPreparationInstrumentJSONRequestBody{}
	c.copyType(&body, validPreparationInstrument)

	res, err := c.authedGeneratedClient.UpdateValidPreparationInstrument(ctx, validPreparationInstrument.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparationInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation update response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidPreparationInstrument archives a valid ingredient preparation.
func (c *Client) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	res, err := c.authedGeneratedClient.ArchiveValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparationInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient preparation archive response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

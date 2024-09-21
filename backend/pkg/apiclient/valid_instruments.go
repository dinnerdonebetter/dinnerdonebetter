package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidInstrument gets a valid instrument.
func (c *Client) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	res, err := c.authedGeneratedClient.GetValidInstrument(ctx, validInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get valid instrument")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid instrument")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRandomValidInstrument gets a valid instrument.
func (c *Client) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	res, err := c.authedGeneratedClient.GetRandomValidInstrument(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get valid instrument")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid instrument")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// SearchValidInstruments searches through a list of valid instruments.
func (c *Client) SearchValidInstruments(ctx context.Context, query string, limit uint8) ([]*types.ValidInstrument, error) {
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

	params := &generated.SearchForValidInstrumentsParams{
		Q:     query,
		Limit: int(limit),
	}

	res, err := c.authedGeneratedClient.SearchForValidInstruments(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "search for valid instruments")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid instruments")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidInstruments retrieves a list of valid instruments.
func (c *Client) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidInstrument], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidInstrumentsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidInstruments(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "valid instruments list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid instruments")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.ValidInstrument]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

// CreateValidInstrument creates a valid instrument.
func (c *Client) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (*types.ValidInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidInstrumentJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidInstrument(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create valid instrument")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid instrument")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidInstrument updates a valid instrument.
func (c *Client) UpdateValidInstrument(ctx context.Context, validInstrument *types.ValidInstrument) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validInstrument == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrument.ID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrument.ID)

	body := generated.UpdateValidInstrumentJSONRequestBody{}
	c.copyType(&body, validInstrument)

	res, err := c.authedGeneratedClient.UpdateValidInstrument(ctx, validInstrument.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update valid instrument")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid instrument %s", validInstrument.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidInstrument archives a valid instrument.
func (c *Client) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	res, err := c.authedGeneratedClient.ArchiveValidInstrument(ctx, validInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive valid instrument")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidInstrument]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid instrument %s", validInstrumentID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidPreparation gets a valid preparation.
func (c *Client) GetValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	res, err := c.authedGeneratedClient.GetValidPreparation(ctx, validPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading valid preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRandomValidPreparation gets a valid preparation.
func (c *Client) GetRandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	res, err := c.authedGeneratedClient.GetRandomValidPreparation(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting random valid preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading random valid preparation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// SearchValidPreparations searches through a list of valid preparations.
func (c *Client) SearchValidPreparations(ctx context.Context, query string, limit uint8) ([]*types.ValidPreparation, error) {
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

	params := &generated.SearchForValidPreparationsParams{
		Q:     query,
		Limit: int(limit),
	}

	res, err := c.authedGeneratedClient.SearchForValidPreparations(ctx, params, c.queryFilterCleaner)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid preparations")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparations")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidPreparations retrieves a list of valid preparations.
func (c *Client) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidPreparationsParams{}
	c.copyType(&params, filter)

	res, err := c.authedGeneratedClient.GetValidPreparations(ctx, params, c.queryFilterCleaner)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid preparations")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparations")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	validPreparations := &types.QueryFilteredResult[types.ValidPreparation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return validPreparations, nil
}

// CreateValidPreparation creates a valid preparation.
func (c *Client) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationRequestInput) (*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidPreparationJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidPreparation(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid preparation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading valid preparation creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidPreparation updates a valid preparation.
func (c *Client) UpdateValidPreparation(ctx context.Context, validPreparation *types.ValidPreparation) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparation == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparation.ID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparation.ID)

	input := generated.UpdateValidPreparationJSONRequestBody{}
	c.copyType(&input, validPreparation)
	res, err := c.authedGeneratedClient.UpdateValidPreparation(ctx, validPreparation.ID, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation %s", validPreparation.ID)
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation %s", validPreparation.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidPreparation archives a valid preparation.
func (c *Client) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	res, err := c.authedGeneratedClient.ArchiveValidPreparation(ctx, validPreparationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid preparation %s", validPreparationID)
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading valid preparation archive response %s", validPreparationID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

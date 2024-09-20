package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// var (\w+) types\.APIResponse\[\*types\.(\w+)\]\n\	if err \= c\.fetchAndUnmarshal\(ctx\, req\, \&(\w+))\)\; err \!\= nil \{

// GetValidVessel gets a valid vessel.
func (c *Client) GetValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	res, err := c.authedGeneratedClient.GetValidVessel(ctx, validVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse types.APIResponse[*types.ValidVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading valid vessel response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRandomValidVessel gets a valid vessel.
func (c *Client) GetRandomValidVessel(ctx context.Context) (*types.ValidVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	res, err := c.authedGeneratedClient.GetRandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse types.APIResponse[*types.ValidVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessel")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// SearchValidVessels searches through a list of valid vessels.
// TODO: add queryFilter param here.
func (c *Client) SearchValidVessels(ctx context.Context, query string, limit uint8) ([]*types.ValidVessel, error) {
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

	// TODO: actually get search query in here
	params := &generated.SearchForValidVesselsParams{
		Q:     query,
		Limit: int(limit),
	}
	res, err := c.authedGeneratedClient.SearchForValidVessels(ctx, params, c.queryFilterCleaner)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse types.APIResponse[[]*types.ValidVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessels")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidVessels retrieves a list of valid vessels.
func (c *Client) GetValidVessels(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidVessel], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidVesselsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidVessels(ctx, params, c.queryFilterCleaner)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessels")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse types.APIResponse[[]*types.ValidVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading valid vessels list response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.ValidVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return output, nil
}

// CreateValidVessel creates a valid vessel.
func (c *Client) CreateValidVessel(ctx context.Context, input *types.ValidVesselCreationRequestInput) (*types.ValidVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidVesselJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidVessel(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse types.APIResponse[*types.ValidVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading valid vessel creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidVessel updates a valid vessel.
func (c *Client) UpdateValidVessel(ctx context.Context, validVessel *types.ValidVessel) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validVessel == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVessel.ID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVessel.ID)

	body := generated.UpdateValidVesselJSONRequestBody{}
	c.copyType(&body, validVessel)

	res, err := c.authedGeneratedClient.UpdateValidVessel(ctx, validVessel.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse types.APIResponse[*types.ValidVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading valid vessel update response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidVessel archives a valid vessel.
func (c *Client) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	res, err := c.authedGeneratedClient.ArchiveValidVessel(ctx, validVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse types.APIResponse[*types.ValidVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading valid vessel archive response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

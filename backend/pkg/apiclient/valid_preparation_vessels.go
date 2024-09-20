package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidPreparationVessel gets a valid preparation vessel.
func (c *Client) GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)

	res, err := c.authedGeneratedClient.GetValidPreparationVessel(ctx, validPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid preparation vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparationVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading valid preparation vessel response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidPreparationVessels retrieves a list of valid preparation vessels.
func (c *Client) GetValidPreparationVessels(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidPreparationVesselsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidPreparationVessels(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid preparation vessels")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidPreparationVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation vessels")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	validPreparationVessels := &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return validPreparationVessels, nil
}

// GetValidPreparationVesselsForPreparation retrieves a list of valid preparation vessels.
func (c *Client) GetValidPreparationVesselsForPreparation(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
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

	params := &generated.GetValidPreparationVesselsByPreparationParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidPreparationVesselsByPreparation(ctx, validPreparationID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation vessels list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidPreparationVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading valid preparation vessels list")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	validPreparationVessels := &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return validPreparationVessels, nil
}

// GetValidPreparationVesselsForVessel retrieves a list of valid preparation vessels.
func (c *Client) GetValidPreparationVesselsForVessel(ctx context.Context, validInstrumentID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
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

	params := &generated.GetValidPreparationVesselsByVesselParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidPreparationVesselsByVessel(ctx, validInstrumentID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation vessels list request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidPreparationVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation vessels")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	validPreparationVessels := &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return validPreparationVessels, nil
}

// CreateValidPreparationVessel creates a valid preparation vessel.
func (c *Client) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselCreationRequestInput) (*types.ValidPreparationVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidPreparationVesselJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidPreparationVessel(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid preparation vessel request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparationVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid preparation vessel")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidPreparationVessel updates a valid preparation vessel.
func (c *Client) UpdateValidPreparationVessel(ctx context.Context, validPreparationVessel *types.ValidPreparationVessel) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVessel == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVessel.ID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVessel.ID)

	body := generated.UpdateValidPreparationVesselJSONRequestBody{}
	c.copyType(&body, validPreparationVessel)

	res, err := c.authedGeneratedClient.UpdateValidPreparationVessel(ctx, validPreparationVessel.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparationVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "parsing valid preparation vessel update response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidPreparationVessel archives a valid preparation vessel.
func (c *Client) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)

	res, err := c.authedGeneratedClient.ArchiveValidPreparationVessel(ctx, validPreparationVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving archive valid preparation vessel")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidPreparationVessel]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "reading valid preparation vessel archive response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetServiceSetting gets a service setting.
func (c *Client) GetServiceSetting(ctx context.Context, serviceSettingID string) (*types.ServiceSetting, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if serviceSettingID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachToSpan(span, keys.ServiceSettingIDKey, serviceSettingID)

	res, err := c.authedGeneratedClient.GetServiceSetting(ctx, serviceSettingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get service setting request")
	}

	var apiResponse *types.APIResponse[*types.ServiceSetting]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service setting")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// SearchServiceSettings searches through a list of service settings.
func (c *Client) SearchServiceSettings(ctx context.Context, query string, limit uint8) ([]*types.ServiceSetting, error) {
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

	params := &generated.SearchForServiceSettingsParams{
		Q:     query,
		Limit: int(limit),
	}

	res, err := c.authedGeneratedClient.SearchForServiceSettings(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building search for service settings request")
	}

	var apiResponse *types.APIResponse[[]*types.ServiceSetting]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetServiceSettings retrieves a list of service settings.
func (c *Client) GetServiceSettings(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSetting], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetServiceSettingsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetServiceSettings(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building service settings list request")
	}

	var apiResponse *types.APIResponse[[]*types.ServiceSetting]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ServiceSetting]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateServiceSetting creates a service setting.
func (c *Client) CreateServiceSetting(ctx context.Context, input *types.ServiceSettingCreationRequestInput) (*types.ServiceSetting, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateServiceSettingJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateServiceSetting(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create service setting request")
	}

	var apiResponse *types.APIResponse[*types.ServiceSetting]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating service setting")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// ArchiveServiceSetting archives a service setting.
func (c *Client) ArchiveServiceSetting(ctx context.Context, serviceSettingID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if serviceSettingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachToSpan(span, keys.ServiceSettingIDKey, serviceSettingID)

	res, err := c.authedGeneratedClient.ArchiveServiceSetting(ctx, serviceSettingID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive service setting request")
	}

	var apiResponse *types.APIResponse[*types.ServiceSetting]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving service setting %s", serviceSettingID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

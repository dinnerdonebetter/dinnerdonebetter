package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetServiceSettingConfigurationForUserByName retrieves a list of service settings.
func (c *Client) GetServiceSettingConfigurationForUserByName(ctx context.Context, settingName string, filter *types.QueryFilter) (*types.ServiceSettingConfiguration, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := c.logger.WithValue(keys.ServiceSettingNameKey, settingName)
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)
	tracing.AttachToSpan(span, keys.ServiceSettingNameKey, settingName)

	params := &generated.GetServiceSettingConfigurationByNameParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetServiceSettingConfigurationByName(ctx, settingName, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "service settings list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ServiceSettingConfiguration]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetServiceSettingConfigurationsForUser retrieves a list of service settings.
func (c *Client) GetServiceSettingConfigurationsForUser(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetServiceSettingConfigurationsForUserParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetServiceSettingConfigurationsForUser(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "service settings list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ServiceSettingConfiguration]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetServiceSettingConfigurationsForHousehold retrieves a list of service settings.
func (c *Client) GetServiceSettingConfigurationsForHousehold(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetServiceSettingConfigurationsForHouseholdParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetServiceSettingConfigurationsForHousehold(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "service settings list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ServiceSettingConfiguration]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateServiceSettingConfiguration creates a service setting.
func (c *Client) CreateServiceSettingConfiguration(ctx context.Context, input *types.ServiceSettingConfigurationCreationRequestInput) (*types.ServiceSettingConfiguration, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateServiceSettingConfigurationJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateServiceSettingConfiguration(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create service setting")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ServiceSettingConfiguration]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating service setting")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateServiceSettingConfiguration updates a service setting.
func (c *Client) UpdateServiceSettingConfiguration(ctx context.Context, serviceSettingConfiguration *types.ServiceSettingConfiguration) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if serviceSettingConfiguration == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ServiceSettingConfigurationIDKey, serviceSettingConfiguration.ID)
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationIDKey, serviceSettingConfiguration.ID)

	body := generated.UpdateServiceSettingConfigurationJSONRequestBody{}
	c.copyType(&body, serviceSettingConfiguration)

	res, err := c.authedGeneratedClient.UpdateServiceSettingConfiguration(ctx, serviceSettingConfiguration.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update service setting")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ServiceSettingConfiguration]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating service setting %s", serviceSettingConfiguration.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveServiceSettingConfiguration archives a service setting.
func (c *Client) ArchiveServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if serviceSettingConfigurationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)

	res, err := c.authedGeneratedClient.ArchiveServiceSettingConfiguration(ctx, serviceSettingConfigurationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive service setting")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ServiceSettingConfiguration]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving service setting %s", serviceSettingConfigurationID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

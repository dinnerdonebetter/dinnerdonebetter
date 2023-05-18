package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetServiceSettingConfigurationForUserByName retrieves a list of service settings.
func (c *Client) GetServiceSettingConfigurationForUserByName(ctx context.Context, settingName string, filter *types.QueryFilter) (*types.ServiceSettingConfiguration, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger).WithValue(keys.ServiceSettingNameKey, settingName)
	tracing.AttachQueryFilterToSpan(span, filter)
	tracing.AttachServiceSettingNameToSpan(span, settingName)

	req, err := c.requestBuilder.BuildGetServiceSettingConfigurationForUserByNameRequest(ctx, settingName, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building service settings list request")
	}

	var serviceSettingConfigurations *types.ServiceSettingConfiguration
	if err = c.fetchAndUnmarshal(ctx, req, &serviceSettingConfigurations); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	return serviceSettingConfigurations, nil
}

// GetServiceSettingConfigurationsForUser retrieves a list of service settings.
func (c *Client) GetServiceSettingConfigurationsForUser(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetServiceSettingConfigurationsForUserRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building service settings list request")
	}

	var serviceSettingConfigurations *types.QueryFilteredResult[types.ServiceSettingConfiguration]
	if err = c.fetchAndUnmarshal(ctx, req, &serviceSettingConfigurations); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	return serviceSettingConfigurations, nil
}

// GetServiceSettingConfigurationsForHousehold retrieves a list of service settings.
func (c *Client) GetServiceSettingConfigurationsForHousehold(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetServiceSettingConfigurationsForHouseholdRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building service settings list request")
	}

	var serviceSettingConfigurations *types.QueryFilteredResult[types.ServiceSettingConfiguration]
	if err = c.fetchAndUnmarshal(ctx, req, &serviceSettingConfigurations); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	return serviceSettingConfigurations, nil
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

	req, err := c.requestBuilder.BuildCreateServiceSettingConfigurationRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create service setting request")
	}

	var serviceSettingConfiguration *types.ServiceSettingConfiguration
	if err = c.fetchAndUnmarshal(ctx, req, &serviceSettingConfiguration); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating service setting")
	}

	return serviceSettingConfiguration, nil
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
	tracing.AttachServiceSettingConfigurationIDToSpan(span, serviceSettingConfiguration.ID)

	req, err := c.requestBuilder.BuildUpdateServiceSettingConfigurationRequest(ctx, serviceSettingConfiguration)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update service setting request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &serviceSettingConfiguration); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating service setting %s", serviceSettingConfiguration.ID)
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
	tracing.AttachServiceSettingConfigurationIDToSpan(span, serviceSettingConfigurationID)

	req, err := c.requestBuilder.BuildArchiveServiceSettingConfigurationRequest(ctx, serviceSettingConfigurationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive service setting request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving service setting %s", serviceSettingConfigurationID)
	}

	return nil
}

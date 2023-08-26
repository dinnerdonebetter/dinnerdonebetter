package apiclient

import (
	"context"

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
	tracing.AttachServiceSettingIDToSpan(span, serviceSettingID)

	req, err := c.requestBuilder.BuildGetServiceSettingRequest(ctx, serviceSettingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get service setting request")
	}

	var serviceSetting *types.ServiceSetting
	if err = c.fetchAndUnmarshal(ctx, req, &serviceSetting); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service setting")
	}

	return serviceSetting, nil
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
		limit = types.DefaultLimit
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	req, err := c.requestBuilder.BuildSearchServiceSettingsRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building search for service settings request")
	}

	var serviceSettings []*types.ServiceSetting
	if err = c.fetchAndUnmarshal(ctx, req, &serviceSettings); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	return serviceSettings, nil
}

// GetServiceSettings retrieves a list of service settings.
func (c *Client) GetServiceSettings(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSetting], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetServiceSettingsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building service settings list request")
	}

	var serviceSettings *types.QueryFilteredResult[types.ServiceSetting]
	if err = c.fetchAndUnmarshal(ctx, req, &serviceSettings); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving service settings")
	}

	return serviceSettings, nil
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

	req, err := c.requestBuilder.BuildCreateServiceSettingRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create service setting request")
	}

	var serviceSetting *types.ServiceSetting
	if err = c.fetchAndUnmarshal(ctx, req, &serviceSetting); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating service setting")
	}

	return serviceSetting, nil
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
	tracing.AttachServiceSettingIDToSpan(span, serviceSettingID)

	req, err := c.requestBuilder.BuildArchiveServiceSettingRequest(ctx, serviceSettingID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive service setting request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving service setting %s", serviceSettingID)
	}

	return nil
}

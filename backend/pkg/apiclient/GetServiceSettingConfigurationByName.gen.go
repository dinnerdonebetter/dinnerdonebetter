// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) GetServiceSettingConfigurationByName(
	ctx context.Context,
	serviceSettingConfigurationName string,
	filter *types.QueryFilter,
) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if serviceSettingConfigurationName == "" {
		return nil, buildInvalidIDError("serviceSettingConfigurationName")
	}
	logger = logger.WithValue(keys.ServiceSettingConfigurationNameKey, serviceSettingConfigurationName)
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationNameKey, serviceSettingConfigurationName)

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, fmt.Sprintf("/api/v1/settings/configurations/user/%s", serviceSettingConfigurationName))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of ServiceSettingConfiguration")
	}

	var apiResponse *types.APIResponse[[]*types.ServiceSettingConfiguration]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of ServiceSettingConfiguration")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

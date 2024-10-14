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

func (c *Client) ArchiveServiceSettingConfiguration(
	ctx context.Context,
	serviceSettingConfigurationID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if serviceSettingConfigurationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/settings/configurations/%s", serviceSettingConfigurationID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a ServiceSettingConfiguration")
	}

	var apiResponse *types.APIResponse[*types.ServiceSettingConfiguration]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading ServiceSettingConfiguration creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

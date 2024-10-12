// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient




import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
)


func (c *Client) ArchiveServiceSetting(
	ctx context.Context,
serviceSettingID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if serviceSettingID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachToSpan(span, keys.ServiceSettingIDKey, serviceSettingID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/settings/%s" , serviceSettingID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a ServiceSetting")
	}

	var apiResponse *types.APIResponse[ *types.ServiceSetting]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading ServiceSetting creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}

	return  nil
}
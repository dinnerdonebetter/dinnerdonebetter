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


func (c *Client) GetUserNotification(
	ctx context.Context,
userNotificationID string,
) ( *types.UserNotification, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if userNotificationID == "" {
		return nil, buildInvalidIDError("userNotification")
	} 
	logger = logger.WithValue(keys.UserNotificationIDKey, userNotificationID)
	tracing.AttachToSpan(span, keys.UserNotificationIDKey, userNotificationID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/user_notifications/%s" , userNotificationID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a UserNotification")
	}

	var apiResponse *types.APIResponse[  *types.UserNotification]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading UserNotification response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}


	return apiResponse.Data, nil
}
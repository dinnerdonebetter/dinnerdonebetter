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

func (c *Client) UpdateUserNotification(
	ctx context.Context,
	userNotificationID string,
	input *types.UserNotificationUpdateRequestInput,
	reqMods ...RequestModifier,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if userNotificationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserNotificationIDKey, userNotificationID)
	tracing.AttachToSpan(span, keys.UserNotificationIDKey, userNotificationID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/user_notifications/%s", userNotificationID))
	req, err := c.buildDataRequest(ctx, http.MethodPatch, u, input)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a UserNotification")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[*types.UserNotification]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading UserNotification creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
)

func (c *Client) GetUserNotifications(
	ctx context.Context,
	filter *QueryFilter,
	reqMods ...RequestModifier,
) (*QueryFilteredResult[UserNotification], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = DefaultQueryFilter()
	}
	// tracing.AttachQueryFilterToSpan(span, filter)

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, "/api/v1/user_notifications")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of UserNotification")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[[]*UserNotification]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of UserNotification")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &QueryFilteredResult[UserNotification]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

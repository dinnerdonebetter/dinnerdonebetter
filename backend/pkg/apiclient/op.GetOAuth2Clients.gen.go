// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
)

func (c *Client) GetOAuth2Clients(
	ctx context.Context,
	filter *QueryFilter,
	reqMods ...RequestModifier,
) (*QueryFilteredResult[OAuth2Client], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = DefaultQueryFilter()
	}
	// tracing.AttachQueryFilterToSpan(span, filter)

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, "/api/v1/oauth2_clients")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of OAuth2Client")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[[]*OAuth2Client]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of OAuth2Client")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &QueryFilteredResult[OAuth2Client]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

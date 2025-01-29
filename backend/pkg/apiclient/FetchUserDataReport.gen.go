// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) FetchUserDataReport(
	ctx context.Context,
	userDataAggregationReportID string,
	reqMods ...RequestModifier,
) (*types.UserDataCollection, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if userDataAggregationReportID == "" {
		return nil, buildInvalidIDError("userDataAggregationReport")
	}
	logger = logger.WithValue(keys.UserDataAggregationReportIDKey, userDataAggregationReportID)
	tracing.AttachToSpan(span, keys.UserDataAggregationReportIDKey, userDataAggregationReportID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/data_privacy/reports/%s", userDataAggregationReportID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a UserDataCollection")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *types.APIResponse[*types.UserDataCollection]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading UserDataCollection response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

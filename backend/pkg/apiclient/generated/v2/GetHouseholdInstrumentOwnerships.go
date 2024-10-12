// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient




import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)


func (c *Client) GetHouseholdInstrumentOwnerships(
	ctx context.Context,
	filter *types.QueryFilter,
) (*types.QueryFilteredResult[types.HouseholdInstrumentOwnership], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

 

	u := c.BuildURL(ctx, filter.ToValues(), fmt.Sprintf("/api/v1/households/instruments" ))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of HouseholdInstrumentOwnership")
	}
	
	var apiResponse *types.APIResponse[ []*types.HouseholdInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of HouseholdInstrumentOwnership")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.HouseholdInstrumentOwnership]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}
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


func (c *Client) ArchiveHouseholdInstrumentOwnership(
	ctx context.Context,
householdInstrumentOwnershipID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/households/instruments/%s" , householdInstrumentOwnershipID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a HouseholdInstrumentOwnership")
	}

	var apiResponse *types.APIResponse[ *types.HouseholdInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading HouseholdInstrumentOwnership creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}

	return  nil
}
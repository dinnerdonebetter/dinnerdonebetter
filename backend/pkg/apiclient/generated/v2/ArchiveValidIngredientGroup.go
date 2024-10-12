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


func (c *Client) ArchiveValidIngredientGroup(
	ctx context.Context,
validIngredientGroupID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientGroupID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_ingredient_groups/%s" , validIngredientGroupID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a ValidIngredientGroup")
	}

	var apiResponse *types.APIResponse[ *types.ValidIngredientGroup]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading ValidIngredientGroup creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}

	return  nil
}
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


func (c *Client) UpdateValidIngredientPreparation(
	ctx context.Context,
validIngredientPreparationID string,
input *types.ValidIngredientPreparationUpdateRequestInput,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientPreparationID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

 


	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_ingredient_preparations/%s" , validIngredientPreparationID ))
	req, err := c.buildDataRequest(ctx, http.MethodPut, u, input)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a ValidIngredientPreparation")
	}

	var apiResponse *types.APIResponse[ *types.ValidIngredientPreparation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading ValidIngredientPreparation creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}


	return nil
}
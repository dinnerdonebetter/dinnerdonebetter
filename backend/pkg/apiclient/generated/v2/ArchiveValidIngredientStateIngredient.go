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


func (c *Client) ArchiveValidIngredientStateIngredient(
	ctx context.Context,
validIngredientStateIngredientID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/valid_ingredient_state_ingredients/%s" , validIngredientStateIngredientID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a ValidIngredientStateIngredient")
	}

	var apiResponse *types.APIResponse[ *types.ValidIngredientStateIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading ValidIngredientStateIngredient creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}

	return  nil
}
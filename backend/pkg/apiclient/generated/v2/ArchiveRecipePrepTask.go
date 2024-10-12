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


func (c *Client) ArchiveRecipePrepTask(
	ctx context.Context,
recipeID string,
recipePrepTaskID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipePrepTaskID == "" {
		return  ErrInvalidIDProvided
	} 
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/recipes/%s/prep_tasks/%s" , recipeID , recipePrepTaskID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "building request to create a RecipePrepTask")
	}

	var apiResponse *types.APIResponse[ *types.RecipePrepTask]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return  observability.PrepareAndLogError(err, logger, span, "loading RecipePrepTask creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return  err
	}

	return  nil
}
// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

func (c *Client) GetRecipePrepTasks(
	ctx context.Context,
	recipeID string,
	filter *QueryFilter,
	reqMods ...RequestModifier,
) (*QueryFilteredResult[RecipePrepTask], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = DefaultQueryFilter()
	}
	// tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	values := filter.ToValues()

	u := c.BuildURL(ctx, values, fmt.Sprintf("/api/v1/recipes/%s/prep_tasks", recipeID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch list of RecipePrepTask")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[[]*RecipePrepTask]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading response for list of RecipePrepTask")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &QueryFilteredResult[RecipePrepTask]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

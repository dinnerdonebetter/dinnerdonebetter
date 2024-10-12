// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"

	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (c *Client) ArchiveUserIngredientPreference(
	ctx context.Context,
	userIngredientPreferenceID string,
) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if userIngredientPreferenceID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)
	tracing.AttachToSpan(span, keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/user_ingredient_preferences/%s", userIngredientPreferenceID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building request to create a UserIngredientPreference")
	}

	var apiResponse *types.APIResponse[*types.UserIngredientPreference]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading UserIngredientPreference creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

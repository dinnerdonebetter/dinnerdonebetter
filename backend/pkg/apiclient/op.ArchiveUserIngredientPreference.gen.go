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

func (c *Client) ArchiveUserIngredientPreference(
	ctx context.Context,
	userIngredientPreferenceID string,
	reqMods ...RequestModifier,
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

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*UserIngredientPreference]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "loading UserIngredientPreference creation response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetUserIngredientPreferences retrieves a list of user ingredient preferences.
func (c *Client) GetUserIngredientPreferences(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.UserIngredientPreference], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetUserIngredientPreferencesParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetUserIngredientPreferences(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building user ingredient preferences list request")
	}

	var apiResponse *types.APIResponse[[]*types.UserIngredientPreference]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving user ingredient preferences")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.UserIngredientPreference]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateUserIngredientPreference creates a user ingredient preference.
func (c *Client) CreateUserIngredientPreference(ctx context.Context, input *types.UserIngredientPreferenceCreationRequestInput) ([]*types.UserIngredientPreference, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateUserIngredientPreferenceJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateUserIngredientPreference(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create user ingredient preference request")
	}

	var apiResponse *types.APIResponse[[]*types.UserIngredientPreference]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating user ingredient preference")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateUserIngredientPreference updates a user ingredient preference.
func (c *Client) UpdateUserIngredientPreference(ctx context.Context, userIngredientPreference *types.UserIngredientPreference) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if userIngredientPreference == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, userIngredientPreference.ID)
	tracing.AttachToSpan(span, keys.UserIngredientPreferenceIDKey, userIngredientPreference.ID)

	body := generated.UpdateUserIngredientPreferenceJSONRequestBody{}
	c.copyType(&body, userIngredientPreference)

	res, err := c.authedGeneratedClient.UpdateUserIngredientPreference(ctx, userIngredientPreference.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update user ingredient preference request")
	}

	var apiResponse *types.APIResponse[types.UserIngredientPreference]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user ingredient preference %s", userIngredientPreference.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveUserIngredientPreference archives a user ingredient preference.
func (c *Client) ArchiveUserIngredientPreference(ctx context.Context, userIngredientPreferenceID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if userIngredientPreferenceID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)
	tracing.AttachToSpan(span, keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)

	res, err := c.authedGeneratedClient.ArchiveUserIngredientPreference(ctx, userIngredientPreferenceID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive user ingredient preference request")
	}

	var apiResponse *types.APIResponse[types.UserIngredientPreference]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving user ingredient preference %s", userIngredientPreferenceID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

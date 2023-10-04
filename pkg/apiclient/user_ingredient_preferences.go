package apiclient

import (
	"context"

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
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetUserIngredientPreferencesRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building user ingredient preferences list request")
	}

	var userIngredientPreferences *types.QueryFilteredResult[types.UserIngredientPreference]
	if err = c.fetchAndUnmarshal(ctx, req, &userIngredientPreferences); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving user ingredient preferences")
	}

	return userIngredientPreferences, nil
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

	req, err := c.requestBuilder.BuildCreateUserIngredientPreferenceRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create user ingredient preference request")
	}

	var userIngredientPreferences []*types.UserIngredientPreference
	if err = c.fetchAndUnmarshal(ctx, req, &userIngredientPreferences); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating user ingredient preference")
	}

	return userIngredientPreferences, nil
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

	req, err := c.requestBuilder.BuildUpdateUserIngredientPreferenceRequest(ctx, userIngredientPreference)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update user ingredient preference request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &userIngredientPreference); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user ingredient preference %s", userIngredientPreference.ID)
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

	req, err := c.requestBuilder.BuildArchiveUserIngredientPreferenceRequest(ctx, userIngredientPreferenceID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive user ingredient preference request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving user ingredient preference %s", userIngredientPreferenceID)
	}

	return nil
}

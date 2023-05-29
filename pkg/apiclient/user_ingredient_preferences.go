package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetUserIngredientPreferences retrieves a list of valid preparations.
func (c *Client) GetUserIngredientPreferences(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.UserIngredientPreference], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetUserIngredientPreferencesRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparations list request")
	}

	var validPreparations *types.QueryFilteredResult[types.UserIngredientPreference]
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparations); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparations")
	}

	return validPreparations, nil
}

// CreateUserIngredientPreference creates a valid preparation.
func (c *Client) CreateUserIngredientPreference(ctx context.Context, input *types.UserIngredientPreferenceCreationRequestInput) (*types.UserIngredientPreference, error) {
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
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid preparation request")
	}

	var validPreparation *types.UserIngredientPreference
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparation); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid preparation")
	}

	return validPreparation, nil
}

// UpdateUserIngredientPreference updates a valid preparation.
func (c *Client) UpdateUserIngredientPreference(ctx context.Context, validPreparation *types.UserIngredientPreference) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparation == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, validPreparation.ID)
	tracing.AttachUserIngredientPreferenceIDToSpan(span, validPreparation.ID)

	req, err := c.requestBuilder.BuildUpdateUserIngredientPreferenceRequest(ctx, validPreparation)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validPreparation); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation %s", validPreparation.ID)
	}

	return nil
}

// ArchiveUserIngredientPreference archives a valid preparation.
func (c *Client) ArchiveUserIngredientPreference(ctx context.Context, validPreparationID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, validPreparationID)
	tracing.AttachUserIngredientPreferenceIDToSpan(span, validPreparationID)

	req, err := c.requestBuilder.BuildArchiveUserIngredientPreferenceRequest(ctx, validPreparationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid preparation request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid preparation %s", validPreparationID)
	}

	return nil
}

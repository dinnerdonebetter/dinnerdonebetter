package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidIngredientStateIngredient gets a valid ingredient state ingredient.
func (c *Client) GetValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)

	req, err := c.requestBuilder.BuildGetValidIngredientStateIngredientRequest(ctx, validIngredientStateIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient state ingredient request")
	}

	var validIngredientStateIngredient *types.ValidIngredientStateIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientStateIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	return validIngredientStateIngredient, nil
}

// GetValidIngredientStateIngredients retrieves a list of valid ingredient state ingredient.
func (c *Client) GetValidIngredientStateIngredients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidIngredientStateIngredientsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient state ingredient list request")
	}

	var validIngredientStateIngredients *types.QueryFilteredResult[types.ValidIngredientStateIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientStateIngredients); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	return validIngredientStateIngredients, nil
}

// GetValidIngredientStateIngredientsForIngredient retrieves a list of valid ingredient state ingredient.
func (c *Client) GetValidIngredientStateIngredientsForIngredient(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildGetValidIngredientStateIngredientsForIngredientRequest(ctx, validIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient state ingredient list request")
	}

	var validPreparationInstruments *types.QueryFilteredResult[types.ValidIngredientStateIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	return validPreparationInstruments, nil
}

// GetValidIngredientStateIngredientsForIngredientState retrieves a list of valid ingredient state ingredient.
func (c *Client) GetValidIngredientStateIngredientsForIngredientState(ctx context.Context, ingredientState string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if ingredientState == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, ingredientState)
	tracing.AttachValidIngredientStateIDToSpan(span, ingredientState)

	req, err := c.requestBuilder.BuildGetValidIngredientStateIngredientsForPreparationRequest(ctx, ingredientState, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient state ingredient list request")
	}

	var validPreparationInstruments *types.QueryFilteredResult[types.ValidIngredientStateIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &validPreparationInstruments); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state ingredient")
	}

	return validPreparationInstruments, nil
}

// CreateValidIngredientStateIngredient creates a valid ingredient state ingredient.
func (c *Client) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientCreationRequestInput) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidIngredientStateIngredientRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient state ingredient request")
	}

	var validIngredientStateIngredient *types.ValidIngredientStateIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientStateIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient state ingredient")
	}

	return validIngredientStateIngredient, nil
}

// UpdateValidIngredientStateIngredient updates a valid ingredient state ingredient.
func (c *Client) UpdateValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredient *types.ValidIngredientStateIngredient) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateIngredient == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredient.ID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredient.ID)

	req, err := c.requestBuilder.BuildUpdateValidIngredientStateIngredientRequest(ctx, validIngredientStateIngredient)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient state ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientStateIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state ingredient %s", validIngredientStateIngredient.ID)
	}

	return nil
}

// ArchiveValidIngredientStateIngredient archives a valid ingredient state ingredient.
func (c *Client) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)

	req, err := c.requestBuilder.BuildArchiveValidIngredientStateIngredientRequest(ctx, validIngredientStateIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient state ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient state ingredient %s", validIngredientStateIngredientID)
	}

	return nil
}

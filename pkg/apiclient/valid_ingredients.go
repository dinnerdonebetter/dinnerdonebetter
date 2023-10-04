package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidIngredient gets a valid ingredient.
func (c *Client) GetValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	req, err := c.requestBuilder.BuildGetValidIngredientRequest(ctx, validIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient request")
	}

	var validIngredient *types.ValidIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient")
	}

	return validIngredient, nil
}

// GetRandomValidIngredient gets a valid ingredient.
func (c *Client) GetRandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	req, err := c.requestBuilder.BuildGetRandomValidIngredientRequest(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient request")
	}

	var validIngredient *types.ValidIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient")
	}

	return validIngredient, nil
}

// SearchValidIngredients searches through a list of valid ingredients.
func (c *Client) SearchValidIngredients(ctx context.Context, query string, limit uint8) (*types.QueryFilteredResult[types.ValidIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	if limit == 0 {
		limit = types.DefaultLimit
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	req, err := c.requestBuilder.BuildSearchValidIngredientsRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building search for valid ingredients request")
	}

	var validIngredients *types.QueryFilteredResult[types.ValidIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredients); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredients")
	}

	return validIngredients, nil
}

// GetValidIngredients retrieves a list of valid ingredients.
func (c *Client) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredient], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidIngredientsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredients list request")
	}

	var validIngredients *types.QueryFilteredResult[types.ValidIngredient]
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredients); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredients")
	}

	return validIngredients, nil
}

// CreateValidIngredient creates a valid ingredient.
func (c *Client) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationRequestInput) (*types.ValidIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidIngredientRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient request")
	}

	var validIngredient *types.ValidIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient")
	}

	return validIngredient, nil
}

// UpdateValidIngredient updates a valid ingredient.
func (c *Client) UpdateValidIngredient(ctx context.Context, validIngredient *types.ValidIngredient) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredient == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredient.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredient.ID)

	req, err := c.requestBuilder.BuildUpdateValidIngredientRequest(ctx, validIngredient)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient %s", validIngredient.ID)
	}

	return nil
}

// ArchiveValidIngredient archives a valid ingredient.
func (c *Client) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	req, err := c.requestBuilder.BuildArchiveValidIngredientRequest(ctx, validIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient %s", validIngredientID)
	}

	return nil
}

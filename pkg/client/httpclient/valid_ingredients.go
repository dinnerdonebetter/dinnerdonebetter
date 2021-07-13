package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// ValidIngredientExists retrieves whether a valid ingredient exists.
func (c *Client) ValidIngredientExists(ctx context.Context, validIngredientID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildValidIngredientExistsRequest(ctx, validIngredientID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building valid ingredient existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for valid ingredient #%d", validIngredientID)
	}

	return exists, nil
}

// GetValidIngredient gets a valid ingredient.
func (c *Client) GetValidIngredient(ctx context.Context, validIngredientID uint64) (*types.ValidIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildGetValidIngredientRequest(ctx, validIngredientID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get valid ingredient request")
	}

	var validIngredient *types.ValidIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid ingredient")
	}

	return validIngredient, nil
}

// SearchValidIngredients searches through a list of valid ingredients.
func (c *Client) SearchValidIngredients(ctx context.Context, query string, limit uint8) ([]*types.ValidIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	if limit == 0 {
		limit = 20
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	req, err := c.requestBuilder.BuildSearchValidIngredientsRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building search for valid ingredients request")
	}

	var validIngredients []*types.ValidIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredients); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid ingredients")
	}

	return validIngredients, nil
}

// GetValidIngredients retrieves a list of valid ingredients.
func (c *Client) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidIngredientsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building valid ingredients list request")
	}

	var validIngredients *types.ValidIngredientList
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredients); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving valid ingredients")
	}

	return validIngredients, nil
}

// CreateValidIngredient creates a valid ingredient.
func (c *Client) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationInput) (*types.ValidIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidIngredientRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create valid ingredient request")
	}

	var validIngredient *types.ValidIngredient
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating valid ingredient")
	}

	return validIngredient, nil
}

// UpdateValidIngredient updates a valid ingredient.
func (c *Client) UpdateValidIngredient(ctx context.Context, validIngredient *types.ValidIngredient) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredient == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredient.ID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredient.ID)

	req, err := c.requestBuilder.BuildUpdateValidIngredientRequest(ctx, validIngredient)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update valid ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient #%d", validIngredient.ID)
	}

	return nil
}

// ArchiveValidIngredient archives a valid ingredient.
func (c *Client) ArchiveValidIngredient(ctx context.Context, validIngredientID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildArchiveValidIngredientRequest(ctx, validIngredientID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive valid ingredient request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving valid ingredient #%d", validIngredientID)
	}

	return nil
}

// GetAuditLogForValidIngredient retrieves a list of audit log entries pertaining to a valid ingredient.
func (c *Client) GetAuditLogForValidIngredient(ctx context.Context, validIngredientID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if validIngredientID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildGetAuditLogForValidIngredientRequest(ctx, validIngredientID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for valid ingredient request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}

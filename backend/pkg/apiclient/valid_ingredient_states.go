package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidIngredientState gets a valid ingredient state.
func (c *Client) GetValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	res, err := c.authedGeneratedClient.GetValidIngredientState(ctx, validIngredientStateID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient state")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientState]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient state response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// SearchValidIngredientStates searches through a list of valid ingredient states.
func (c *Client) SearchValidIngredientStates(ctx context.Context, query string, limit uint8) ([]*types.ValidIngredientState, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	if limit == 0 {
		limit = types.DefaultQueryFilterLimit
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	params := &generated.SearchForValidIngredientStatesParams{
		Q:     query,
		Limit: int(limit),
	}

	res, err := c.authedGeneratedClient.SearchForValidIngredientStates(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid ingredient states")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientState]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient states response")
	}

	return apiResponse.Data, nil
}

// GetValidIngredientStates retrieves a list of valid ingredient states.
func (c *Client) GetValidIngredientStates(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientState], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidIngredientStatesParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientStates(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "listing valid ingredient states")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientState]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient states list response")
	}

	response := &types.QueryFilteredResult[types.ValidIngredientState]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateValidIngredientState creates a valid ingredient state.
func (c *Client) CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateCreationRequestInput) (*types.ValidIngredientState, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidIngredientStateJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidIngredientState(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient state")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientState]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient state creation response")
	}

	return apiResponse.Data, nil
}

// UpdateValidIngredientState updates a valid ingredient state.
func (c *Client) UpdateValidIngredientState(ctx context.Context, validIngredientState *types.ValidIngredientState) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientState == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientState.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientState.ID)

	body := generated.UpdateValidIngredientStateJSONRequestBody{}
	c.copyType(&body, validIngredientState)

	res, err := c.authedGeneratedClient.UpdateValidIngredientState(ctx, validIngredientState.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientState]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient state update response")
	}

	return nil
}

// ArchiveValidIngredientState archives a valid ingredient state.
func (c *Client) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	res, err := c.authedGeneratedClient.ArchiveValidIngredientState(ctx, validIngredientStateID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient state")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientState]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "parsing valid ingredient state archive response")
	}

	return nil
}

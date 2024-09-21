package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidIngredientGroup gets a valid ingredient group.
func (c *Client) GetValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*types.ValidIngredientGroup, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientGroupID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	res, err := c.authedGeneratedClient.GetValidIngredientGroup(ctx, validIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get valid ingredient group")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientGroup]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient group")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// SearchValidIngredientGroups searches through a list of valid ingredient groups.
func (c *Client) SearchValidIngredientGroups(ctx context.Context, query string, limit uint8) ([]*types.ValidIngredientGroup, error) {
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

	params := &generated.SearchForValidIngredientGroupsParams{
		Q:     query,
		Limit: int(limit),
	}

	res, err := c.authedGeneratedClient.SearchForValidIngredientGroups(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "search for valid ingredient groups")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientGroup]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient groups")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidIngredientGroups retrieves a list of valid ingredient groups.
func (c *Client) GetValidIngredientGroups(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientGroup], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetValidIngredientGroupsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetValidIngredientGroups(ctx, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "valid ingredients list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.ValidIngredientGroup]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredients")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.ValidIngredientGroup]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateValidIngredientGroup creates a valid ingredient.
func (c *Client) CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupCreationRequestInput) (*types.ValidIngredientGroup, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateValidIngredientGroupJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateValidIngredientGroup(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create valid ingredient")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientGroup]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient group")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidIngredientGroup updates a valid ingredient group.
func (c *Client) UpdateValidIngredientGroup(ctx context.Context, validIngredientGroup *types.ValidIngredientGroup) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientGroup == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroup.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroup.ID)

	body := generated.UpdateValidIngredientGroupJSONRequestBody{}
	c.copyType(&body, validIngredientGroup)

	res, err := c.authedGeneratedClient.UpdateValidIngredientGroup(ctx, validIngredientGroup.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update valid ingredient group")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientGroup]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient group %s", validIngredientGroup.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidIngredientGroup archives a valid ingredient group.
func (c *Client) ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientGroupID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	res, err := c.authedGeneratedClient.ArchiveValidIngredientGroup(ctx, validIngredientGroupID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive valid ingredient group")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.ValidIngredientGroup]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient group %s", validIngredientGroupID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

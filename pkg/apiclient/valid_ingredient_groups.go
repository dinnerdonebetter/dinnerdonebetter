package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
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
	tracing.AttachValidIngredientGroupIDToSpan(span, validIngredientGroupID)

	req, err := c.requestBuilder.BuildGetValidIngredientGroupRequest(ctx, validIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient group request")
	}

	var validIngredientGroup *types.ValidIngredientGroup
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientGroup); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient group")
	}

	return validIngredientGroup, nil
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
		limit = 20
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	req, err := c.requestBuilder.BuildSearchValidIngredientGroupsRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building search for valid ingredient groups request")
	}

	var validIngredientGroups []*types.ValidIngredientGroup
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientGroups); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient groups")
	}

	return validIngredientGroups, nil
}

// GetValidIngredientGroups retrieves a list of valid ingredient groups.
func (c *Client) GetValidIngredientGroups(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientGroup], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidIngredientGroupsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredients list request")
	}

	var validIngredientGroups *types.QueryFilteredResult[types.ValidIngredientGroup]
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientGroups); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredients")
	}

	return validIngredientGroups, nil
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

	req, err := c.requestBuilder.BuildCreateValidIngredientGroupRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient request")
	}

	var validIngredientGroup *types.ValidIngredientGroup
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientGroup); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient group")
	}

	return validIngredientGroup, nil
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
	tracing.AttachValidIngredientGroupIDToSpan(span, validIngredientGroup.ID)

	req, err := c.requestBuilder.BuildUpdateValidIngredientGroupRequest(ctx, validIngredientGroup)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient group request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validIngredientGroup); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient group %s", validIngredientGroup.ID)
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
	tracing.AttachValidIngredientGroupIDToSpan(span, validIngredientGroupID)

	req, err := c.requestBuilder.BuildArchiveValidIngredientGroupRequest(ctx, validIngredientGroupID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient group request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient group %s", validIngredientGroupID)
	}

	return nil
}

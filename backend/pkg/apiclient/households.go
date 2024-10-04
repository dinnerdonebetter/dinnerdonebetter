package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetCurrentHousehold retrieves a household.
func (c *Client) GetCurrentHousehold(ctx context.Context) (*types.Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	res, err := c.authedGeneratedClient.GetActiveHousehold(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "household retrieval")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetHousehold retrieves a household.
func (c *Client) GetHousehold(ctx context.Context, householdID string) (*types.Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	res, err := c.authedGeneratedClient.GetHousehold(ctx, householdID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "household retrieval")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetHouseholds retrieves a list of households.
func (c *Client) GetHouseholds(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Household], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetHouseholdsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetHouseholds(ctx, params)
	if err != nil {
		return nil, observability.PrepareError(err, span, "household list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving households")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.Household]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

// CreateHousehold creates a household.
func (c *Client) CreateHousehold(ctx context.Context, input *types.HouseholdCreationRequestInput) (*types.Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	body := generated.CreateHouseholdJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateHousehold(ctx, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "household creation")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "creating household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateHousehold updates a household.
func (c *Client) UpdateHousehold(ctx context.Context, household *types.Household) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if household == nil {
		return ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, household.ID)

	body := generated.UpdateHouseholdJSONRequestBody{}
	c.copyType(&body, household)

	res, err := c.authedGeneratedClient.UpdateHousehold(ctx, household.ID, body)
	if err != nil {
		return observability.PrepareError(err, span, "household update")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "updating household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveHousehold archives a household.
func (c *Client) ArchiveHousehold(ctx context.Context, householdID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	res, err := c.authedGeneratedClient.ArchiveHousehold(ctx, householdID)
	if err != nil {
		return observability.PrepareError(err, span, "household archive")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "archiving household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// InviteUserToHousehold adds a user to a household.
func (c *Client) InviteUserToHousehold(ctx context.Context, destinationHouseholdID string, input *types.HouseholdInvitationCreationRequestInput) (*types.HouseholdInvitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, destinationHouseholdID)

	// we don't validate here because it needs to have the user ID

	body := generated.CreateHouseholdInvitationJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateHouseholdInvitation(ctx, destinationHouseholdID, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "add user to household")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "adding user to household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// MarkAsDefault marks a given household as the default for a given user.
func (c *Client) MarkAsDefault(ctx context.Context, householdID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	res, err := c.authedGeneratedClient.SetDefaultHousehold(ctx, householdID)
	if err != nil {
		return observability.PrepareError(err, span, "mark household as default")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "marking household as default")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// RemoveUserFromHousehold removes a user from a household.
func (c *Client) RemoveUserFromHousehold(ctx context.Context, householdID, userID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	res, err := c.authedGeneratedClient.ArchiveUserMembership(ctx, householdID, userID)
	if err != nil {
		return observability.PrepareError(err, span, "remove user from household")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "removing user from household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ModifyMemberPermissions modifies a given user's permissions for a given household.
func (c *Client) ModifyMemberPermissions(ctx context.Context, householdID, userID string, input *types.ModifyUserPermissionsInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "validating input")
	}

	body := generated.UpdateHouseholdMemberPermissionsJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.UpdateHouseholdMemberPermissions(ctx, householdID, userID, body)
	if err != nil {
		return observability.PrepareError(err, span, "modify household member permissions")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "modifying user household permissions")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// TransferHouseholdOwnership transfers ownership of a household to a given user.
func (c *Client) TransferHouseholdOwnership(ctx context.Context, householdID string, input *types.HouseholdOwnershipTransferInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	tracing.AttachToSpan(span, "old_owner", input.CurrentOwner)
	tracing.AttachToSpan(span, "new_owner", input.NewOwner)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "validating input")
	}

	body := generated.TransferHouseholdOwnershipJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.TransferHouseholdOwnership(ctx, householdID, body)
	if err != nil {
		return observability.PrepareError(err, span, "transfer household ownership")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "transferring household to user")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

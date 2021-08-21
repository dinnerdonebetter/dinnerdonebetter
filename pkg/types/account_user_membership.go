package types

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// HouseholdUserMembership defines a relationship between a user and an household.
	HouseholdUserMembership struct {
		ArchivedOn         *uint64  `json:"archivedOn"`
		LastUpdatedOn      *uint64  `json:"lastUpdatedOn"`
		HouseholdRoles     []string `json:"householdRole"`
		BelongsToUser      uint64   `json:"belongsToUser"`
		BelongsToHousehold uint64   `json:"belongsToHousehold"`
		CreatedOn          uint64   `json:"createdOn"`
		ID                 uint64   `json:"id"`
		DefaultHousehold   bool     `json:"defaultHousehold"`
	}

	// HouseholdUserMembershipList represents a list of household user memberships.
	HouseholdUserMembershipList struct {
		HouseholdUserMemberships []*HouseholdUserMembership `json:"householdUserMemberships"`
		Pagination
	}

	// HouseholdUserMembershipCreationInput represents what a User could set as input for creating household user memberships.
	HouseholdUserMembershipCreationInput struct {
		BelongsToUser      uint64 `json:"belongsToUser"`
		BelongsToHousehold uint64 `json:"belongsToHousehold"`
	}

	// HouseholdUserMembershipUpdateInput represents what a User could set as input for updating household user memberships.
	HouseholdUserMembershipUpdateInput struct {
		BelongsToUser      uint64 `json:"belongsToUser"`
		BelongsToHousehold uint64 `json:"belongsToHousehold"`
	}

	// AddUserToHouseholdInput represents what a User could set as input for updating household user memberships.
	AddUserToHouseholdInput struct {
		Reason         string   `json:"reason"`
		HouseholdRoles []string `json:"householdRole"`
		UserID         uint64   `json:"userID"`
		HouseholdID    uint64   `json:"householdID"`
	}

	// HouseholdOwnershipTransferInput represents what a User could set as input for updating household user memberships.
	HouseholdOwnershipTransferInput struct {
		Reason       string `json:"reason"`
		CurrentOwner uint64 `json:"currentOwner"`
		NewOwner     uint64 `json:"newOwner"`
	}

	// ModifyUserPermissionsInput  represents what a User could set as input for updating household user memberships.
	ModifyUserPermissionsInput struct {
		Reason   string   `json:"reason"`
		NewRoles []string `json:"newRoles"`
	}

	// HouseholdUserMembershipDataManager describes a structure capable of storing householdUserMemberships permanently.
	HouseholdUserMembershipDataManager interface {
		BuildSessionContextDataForUser(ctx context.Context, userID uint64) (*SessionContextData, error)
		GetDefaultHouseholdIDForUser(ctx context.Context, userID uint64) (uint64, error)
		MarkHouseholdAsUserDefault(ctx context.Context, userID, householdID, changedByUser uint64) error
		UserIsMemberOfHousehold(ctx context.Context, userID, householdID uint64) (bool, error)
		ModifyUserPermissions(ctx context.Context, householdID, userID, changedByUser uint64, input *ModifyUserPermissionsInput) error
		TransferHouseholdOwnership(ctx context.Context, householdID uint64, transferredBy uint64, input *HouseholdOwnershipTransferInput) error
		AddUserToHousehold(ctx context.Context, input *AddUserToHouseholdInput, addedByUser uint64) error
		RemoveUserFromHousehold(ctx context.Context, userID, householdID, removedByUser uint64, reason string) error
	}
)

var _ validation.ValidatableWithContext = (*AddUserToHouseholdInput)(nil)

// ValidateWithContext validates an AddUserToHouseholdInput.
func (x *AddUserToHouseholdInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.UserID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdOwnershipTransferInput)(nil)

// ValidateWithContext validates a HouseholdOwnershipTransferInput.
func (x *HouseholdOwnershipTransferInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.CurrentOwner, validation.Required),
		validation.Field(&x.NewOwner, validation.Required),
		validation.Field(&x.Reason, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ModifyUserPermissionsInput)(nil)

// ValidateWithContext validates a ModifyUserPermissionsInput.
func (x *ModifyUserPermissionsInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.NewRoles, validation.Required),
		validation.Field(&x.Reason, validation.Required),
	)
}

package types

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// UserMembershipDataType indicates an event is user membership-related.
	UserMembershipDataType dataType = "user_membership"
)

type (
	// HouseholdUserMembership defines a relationship between a user and a household.
	HouseholdUserMembership struct {
		_ struct{}

		ArchivedOn         *uint64  `json:"archivedOn"`
		LastUpdatedOn      *uint64  `json:"lastUpdatedOn"`
		ID                 string   `json:"id"`
		BelongsToUser      string   `json:"belongsToUser"`
		BelongsToHousehold string   `json:"belongsToHousehold"`
		HouseholdRoles     []string `json:"householdRole"`
		CreatedOn          uint64   `json:"createdOn"`
		DefaultHousehold   bool     `json:"defaultHousehold"`
	}

	// HouseholdUserMembershipList represents a list of household user memberships.
	HouseholdUserMembershipList struct {
		_ struct{}

		HouseholdUserMemberships []*HouseholdUserMembership `json:"householdUserMemberships"`
		Pagination
	}

	// HouseholdUserMembershipCreationInput represents what a User could set as input for creating household user memberships.
	HouseholdUserMembershipCreationInput struct {
		_ struct{}

		ID                 string `json:"-"`
		BelongsToUser      string `json:"belongsToUser"`
		BelongsToHousehold string `json:"belongsToHousehold"`
	}

	// HouseholdUserMembershipUpdateInput represents what a User could set as input for updating household user memberships.
	HouseholdUserMembershipUpdateInput struct {
		_ struct{}

		BelongsToUser      string `json:"belongsToUser"`
		BelongsToHousehold string `json:"belongsToHousehold"`
	}

	// HouseholdUserMembershipCreationRequestInput represents what a User could set as input for updating household user memberships.
	HouseholdUserMembershipCreationRequestInput struct {
		_ struct{}

		ID             string   `json:"-"`
		Reason         string   `json:"reason"`
		UserID         string   `json:"userID"`
		HouseholdID    string   `json:"-"`
		HouseholdRoles []string `json:"-"`
	}

	// HouseholdUserMembershipDatabaseCreationInput represents what a User could set as input for updating household user memberships.
	HouseholdUserMembershipDatabaseCreationInput struct {
		_ struct{}

		ID             string   `json:"-"`
		Reason         string   `json:"reason"`
		UserID         string   `json:"userID"`
		HouseholdID    string   `json:"householdID"`
		HouseholdRoles []string `json:"householdRoles"`
	}

	// HouseholdOwnershipTransferInput represents what a User could set as input for updating household user memberships.
	HouseholdOwnershipTransferInput struct {
		_ struct{}

		Reason       string `json:"reason"`
		CurrentOwner string `json:"currentOwner"`
		NewOwner     string `json:"newOwner"`
	}

	// ModifyUserPermissionsInput  represents what a User could set as input for updating household user memberships.
	ModifyUserPermissionsInput struct {
		_ struct{}

		Reason   string   `json:"reason"`
		NewRoles []string `json:"newRoles"`
	}

	// HouseholdUserMembershipDataManager describes a structure capable of storing householdUserMemberships permanently.
	HouseholdUserMembershipDataManager interface {
		BuildSessionContextDataForUser(ctx context.Context, userID string) (*SessionContextData, error)
		GetDefaultHouseholdIDForUser(ctx context.Context, userID string) (string, error)
		MarkHouseholdAsUserDefault(ctx context.Context, userID, householdID string) error
		UserIsMemberOfHousehold(ctx context.Context, userID, householdID string) (bool, error)
		ModifyUserPermissions(ctx context.Context, householdID, userID string, input *ModifyUserPermissionsInput) error
		TransferHouseholdOwnership(ctx context.Context, householdID string, input *HouseholdOwnershipTransferInput) error
		RemoveUserFromHousehold(ctx context.Context, userID, householdID string) error
	}
)

var _ validation.ValidatableWithContext = (*HouseholdUserMembershipCreationRequestInput)(nil)

// ValidateWithContext validates an HouseholdUserMembershipCreationRequestInput.
func (x *HouseholdUserMembershipCreationRequestInput) ValidateWithContext(ctx context.Context) error {
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

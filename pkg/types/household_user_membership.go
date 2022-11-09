package types

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// UserMembershipDataType indicates an event is user membership-related.
	UserMembershipDataType dataType = "user_membership"

	// HouseholdUserMembershipCreatedCustomerEventType indicates a household user membership was created.
	HouseholdUserMembershipCreatedCustomerEventType CustomerEventType = "household_user_membership_created"
	// HouseholdUserMembershipUpdatedCustomerEventType indicates a household user membership was updated.
	HouseholdUserMembershipUpdatedCustomerEventType CustomerEventType = "household_user_membership_updated"
	// HouseholdUserMembershipArchivedCustomerEventType indicates a household user membership was archived.
	HouseholdUserMembershipArchivedCustomerEventType CustomerEventType = "household_user_membership_archived"
)

type (
	// HouseholdUserMembership defines a relationship between a user and a household.
	HouseholdUserMembership struct {
		_                  struct{}
		CreatedAt          time.Time  `json:"createdAt"`
		LastUpdatedAt      *time.Time `json:"lastUpdatedAt"`
		ArchivedAt         *time.Time `json:"archivedAt"`
		ID                 string     `json:"id"`
		BelongsToUser      string     `json:"belongsToUser"`
		BelongsToHousehold string     `json:"belongsToHousehold"`
		HouseholdRole      string     `json:"householdRoles"`
		DefaultHousehold   bool       `json:"defaultHousehold"`
	}

	// HouseholdUserMembershipWithUser defines a relationship between a user and a household.
	HouseholdUserMembershipWithUser struct {
		_                  struct{}
		CreatedAt          time.Time  `json:"createdAt"`
		LastUpdatedAt      *time.Time `json:"lastUpdatedAt"`
		BelongsToUser      *User      `json:"belongsToUser"`
		ArchivedAt         *time.Time `json:"archivedAt"`
		ID                 string     `json:"id"`
		BelongsToHousehold string     `json:"belongsToHousehold"`
		HouseholdRole      string     `json:"householdRoles"`
		DefaultHousehold   bool       `json:"defaultHousehold"`
	}

	// HouseholdUserMembershipList represents a list of household user memberships.
	HouseholdUserMembershipList struct {
		_ struct{}

		HouseholdUserMemberships []*HouseholdUserMembership `json:"data"`
		Pagination
	}

	// HouseholdUserMembershipCreationRequestInput represents what a User could set as input for updating household user memberships.
	HouseholdUserMembershipCreationRequestInput struct {
		_ struct{}

		ID            string `json:"-"`
		Reason        string `json:"reason"`
		UserID        string `json:"userID"`
		HouseholdID   string `json:"-"`
		HouseholdRole string `json:"-"`
	}

	// HouseholdUserMembershipDatabaseCreationInput represents what a User could set as input for updating household user memberships.
	HouseholdUserMembershipDatabaseCreationInput struct {
		_ struct{}

		ID            string `json:"-"`
		Reason        string `json:"reason"`
		UserID        string `json:"userID"`
		HouseholdID   string `json:"householdID"`
		HouseholdRole string `json:"householdRoles"`
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

		Reason  string `json:"reason"`
		NewRole string `json:"newRoles"`
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
		validation.Field(&x.NewRole, validation.Required),
		validation.Field(&x.Reason, validation.Required),
	)
}

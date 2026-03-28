package customroles

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authorization"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ServiceScope indicates a custom role that is globally scoped.
	ServiceScope = "service"
	// AccountScope indicates a custom role that is scoped to a specific account.
	AccountScope = "account"

	// CustomRoleCreatedServiceEventType indicates a custom role was created.
	CustomRoleCreatedServiceEventType = "custom_role_created"
	// CustomRoleUpdatedServiceEventType indicates a custom role was updated.
	CustomRoleUpdatedServiceEventType = "custom_role_updated"
	// CustomRoleArchivedServiceEventType indicates a custom role was archived.
	CustomRoleArchivedServiceEventType = "custom_role_archived"
	// CustomRoleAssignedServiceEventType indicates a custom role was assigned.
	CustomRoleAssignedServiceEventType = "custom_role_assigned"
	// CustomRoleUnassignedServiceEventType indicates a custom role was unassigned.
	CustomRoleUnassignedServiceEventType = "custom_role_unassigned"
)

type (
	// CustomRole represents an admin-defined role with a subset of permissions.
	CustomRole struct {
		_ struct{} `json:"-"`

		CreatedAt        time.Time  `json:"createdAt"`
		LastUpdatedAt    *time.Time `json:"lastUpdatedAt"`
		ArchivedAt       *time.Time `json:"archivedAt"`
		ID               string     `json:"id"`
		Name             string     `json:"name"`
		Description      string     `json:"description"`
		Scope            string     `json:"scope"`
		BelongsToAccount string     `json:"belongsToAccount,omitempty"`
		CreatedBy        string     `json:"createdBy"`
		Permissions      []string   `json:"permissions"`
	}

	// CustomRoleCreationRequestInput is the input for creating a custom role.
	CustomRoleCreationRequestInput struct {
		_ struct{} `json:"-"`

		Name             string   `json:"name"`
		Description      string   `json:"description"`
		Scope            string   `json:"scope"`
		BelongsToAccount string   `json:"belongsToAccount,omitempty"`
		Permissions      []string `json:"permissions"`
	}

	// CustomRoleDatabaseCreationInput is the DB-ready input for creating a custom role.
	CustomRoleDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID               string   `json:"-"`
		Name             string   `json:"-"`
		Description      string   `json:"-"`
		Scope            string   `json:"-"`
		BelongsToAccount string   `json:"-"`
		CreatedBy        string   `json:"-"`
		Permissions      []string `json:"-"`
	}

	// CustomRoleUpdateRequestInput is the input for updating a custom role.
	CustomRoleUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name        *string  `json:"name,omitempty"`
		Description *string  `json:"description,omitempty"`
		Permissions []string `json:"permissions,omitempty"`
	}

	// CustomRoleAssignment represents an assignment of a custom role to a user or membership.
	CustomRoleAssignment struct {
		_ struct{} `json:"-"`

		CreatedAt           time.Time  `json:"createdAt"`
		ArchivedAt          *time.Time `json:"archivedAt"`
		ID                  string     `json:"id"`
		CustomRoleID        string     `json:"customRoleID"`
		CustomRoleName      string     `json:"customRoleName"`
		UserID              string     `json:"userID,omitempty"`
		AccountMembershipID string     `json:"accountMembershipID,omitempty"`
	}

	// CustomRoleAssignmentCreationInput is the input for assigning a custom role.
	CustomRoleAssignmentCreationInput struct {
		_ struct{} `json:"-"`

		CustomRoleID        string `json:"customRoleID"`
		UserID              string `json:"userID,omitempty"`
		AccountMembershipID string `json:"accountMembershipID,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*CustomRoleCreationRequestInput)(nil)

// ValidateWithContext validates a CustomRoleCreationRequestInput.
func (x *CustomRoleCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	allPerms := authorization.AllPermissionIDs()

	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Scope, validation.Required, validation.In(ServiceScope, AccountScope)),
		validation.Field(&x.Permissions,
			validation.Required,
			validation.Each(validation.By(func(value any) error {
				s, ok := value.(string)
				if !ok {
					return validation.ErrInInvalid
				}
				if _, exists := allPerms[s]; !exists {
					return validation.NewError("validation_unknown_permission", "unknown permission: "+s)
				}
				return nil
			})),
		),
	)
}

var _ validation.ValidatableWithContext = (*CustomRoleUpdateRequestInput)(nil)

// ValidateWithContext validates a CustomRoleUpdateRequestInput.
func (x *CustomRoleUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	allPerms := authorization.AllPermissionIDs()

	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Permissions,
			validation.Each(validation.By(func(value any) error {
				s, ok := value.(string)
				if !ok {
					return validation.ErrInInvalid
				}
				if _, exists := allPerms[s]; !exists {
					return validation.NewError("validation_unknown_permission", "unknown permission: "+s)
				}
				return nil
			})),
		),
	)
}

var _ validation.ValidatableWithContext = (*CustomRoleAssignmentCreationInput)(nil)

// ValidateWithContext validates a CustomRoleAssignmentCreationInput.
func (x *CustomRoleAssignmentCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.CustomRoleID, validation.Required),
	)
}

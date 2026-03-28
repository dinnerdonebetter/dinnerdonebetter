package customroles

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
)

// Repository describes a structure capable of storing custom roles permanently.
type Repository interface {
	CreateCustomRole(ctx context.Context, input *CustomRoleDatabaseCreationInput) (*CustomRole, error)
	GetCustomRole(ctx context.Context, roleID string) (*CustomRole, error)
	ListServiceScopedCustomRoles(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[CustomRole], error)
	ListAccountScopedCustomRoles(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[CustomRole], error)
	UpdateCustomRole(ctx context.Context, roleID string, input *CustomRoleUpdateRequestInput) error
	ArchiveCustomRole(ctx context.Context, roleID string) error

	AssignCustomRoleToUser(ctx context.Context, id, roleID, userID string) error
	AssignCustomRoleToMembership(ctx context.Context, id, roleID, membershipID string) error
	UnassignCustomRole(ctx context.Context, assignmentID string) error
	GetCustomRoleAssignmentsForUser(ctx context.Context, userID string) ([]*CustomRoleAssignment, error)
	GetCustomRoleAssignmentsForMembership(ctx context.Context, membershipID string) ([]*CustomRoleAssignment, error)

	GetServiceScopedRoleIDsForUser(ctx context.Context, userID string) ([]string, error)
	GetAccountScopedRoleIDsForMembership(ctx context.Context, membershipID string) ([]string, error)

	// GetAllCustomRolePermissions returns all role->permission mappings for cache warmup.
	GetAllCustomRolePermissions(ctx context.Context) (map[string][]string, error)
}

package manager

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/customroles"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
)

// CustomRolesDataManager describes a structure capable of managing custom roles.
type CustomRolesDataManager interface {
	CreateCustomRole(ctx context.Context, userID string, input *customroles.CustomRoleCreationRequestInput) (*customroles.CustomRole, error)
	GetCustomRole(ctx context.Context, roleID string) (*customroles.CustomRole, error)
	ListServiceScopedCustomRoles(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[customroles.CustomRole], error)
	ListAccountScopedCustomRoles(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[customroles.CustomRole], error)
	UpdateCustomRole(ctx context.Context, roleID string, input *customroles.CustomRoleUpdateRequestInput) error
	ArchiveCustomRole(ctx context.Context, roleID string) error

	AssignRoleToUser(ctx context.Context, roleID, userID string) (*customroles.CustomRoleAssignment, error)
	UnassignRoleFromUser(ctx context.Context, assignmentID string) error
	AssignRoleToMembership(ctx context.Context, roleID, membershipID string) (*customroles.CustomRoleAssignment, error)
	UnassignRoleFromMembership(ctx context.Context, assignmentID string) error
	GetRoleAssignmentsForUser(ctx context.Context, userID string) ([]*customroles.CustomRoleAssignment, error)
	GetRoleAssignmentsForMembership(ctx context.Context, membershipID string) ([]*customroles.CustomRoleAssignment, error)

	ListAvailablePermissions(ctx context.Context) []string
}

package manager

import (
	"context"
	"sort"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/customroles"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

var _ CustomRolesDataManager = (*customRolesManager)(nil)

type customRolesManager struct {
	tracer     tracing.Tracer
	logger     logging.Logger
	repository customroles.Repository
	cache      *authorization.RolePermissionCache
}

// NewCustomRolesDataManager creates a new CustomRolesDataManager.
func NewCustomRolesDataManager(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	repository customroles.Repository,
	cache *authorization.RolePermissionCache,
) CustomRolesDataManager {
	return &customRolesManager{
		tracer:     tracing.NewNamedTracer(tracerProvider, "custom_roles_manager"),
		logger:     logging.NewNamedLogger(logger, "custom_roles_manager"),
		repository: repository,
		cache:      cache,
	}
}

func (m *customRolesManager) refreshCache(ctx context.Context) {
	if err := m.cache.Refresh(ctx, m.repository.GetAllCustomRolePermissions); err != nil {
		m.logger.Error("refreshing custom role permission cache", err)
	}
}

func (m *customRolesManager) CreateCustomRole(ctx context.Context, userID string, input *customroles.CustomRoleCreationRequestInput) (*customroles.CustomRole, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	dbInput := &customroles.CustomRoleDatabaseCreationInput{
		ID:               identifiers.New(),
		Name:             input.Name,
		Description:      input.Description,
		Scope:            input.Scope,
		BelongsToAccount: input.BelongsToAccount,
		CreatedBy:        userID,
		Permissions:      input.Permissions,
	}

	role, err := m.repository.CreateCustomRole(ctx, dbInput)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating custom role")
	}

	m.refreshCache(ctx)

	return role, nil
}

func (m *customRolesManager) GetCustomRole(ctx context.Context, roleID string) (*customroles.CustomRole, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repository.GetCustomRole(ctx, roleID)
}

func (m *customRolesManager) ListServiceScopedCustomRoles(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[customroles.CustomRole], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repository.ListServiceScopedCustomRoles(ctx, filter)
}

func (m *customRolesManager) ListAccountScopedCustomRoles(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[customroles.CustomRole], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repository.ListAccountScopedCustomRoles(ctx, accountID, filter)
}

func (m *customRolesManager) UpdateCustomRole(ctx context.Context, roleID string, input *customroles.CustomRoleUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if err := m.repository.UpdateCustomRole(ctx, roleID, input); err != nil {
		return observability.PrepareError(err, span, "updating custom role")
	}

	m.refreshCache(ctx)

	return nil
}

func (m *customRolesManager) ArchiveCustomRole(ctx context.Context, roleID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if err := m.repository.ArchiveCustomRole(ctx, roleID); err != nil {
		return observability.PrepareError(err, span, "archiving custom role")
	}

	m.refreshCache(ctx)

	return nil
}

func (m *customRolesManager) AssignRoleToUser(ctx context.Context, roleID, userID string) (*customroles.CustomRoleAssignment, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	id := identifiers.New()
	if err := m.repository.AssignCustomRoleToUser(ctx, id, roleID, userID); err != nil {
		return nil, observability.PrepareError(err, span, "assigning role to user")
	}

	return &customroles.CustomRoleAssignment{
		ID:           id,
		CustomRoleID: roleID,
		UserID:       userID,
	}, nil
}

func (m *customRolesManager) UnassignRoleFromUser(ctx context.Context, assignmentID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repository.UnassignCustomRole(ctx, assignmentID)
}

func (m *customRolesManager) AssignRoleToMembership(ctx context.Context, roleID, membershipID string) (*customroles.CustomRoleAssignment, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	id := identifiers.New()
	if err := m.repository.AssignCustomRoleToMembership(ctx, id, roleID, membershipID); err != nil {
		return nil, observability.PrepareError(err, span, "assigning role to membership")
	}

	return &customroles.CustomRoleAssignment{
		ID:                  id,
		CustomRoleID:        roleID,
		AccountMembershipID: membershipID,
	}, nil
}

func (m *customRolesManager) UnassignRoleFromMembership(ctx context.Context, assignmentID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repository.UnassignCustomRole(ctx, assignmentID)
}

func (m *customRolesManager) GetRoleAssignmentsForUser(ctx context.Context, userID string) ([]*customroles.CustomRoleAssignment, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repository.GetCustomRoleAssignmentsForUser(ctx, userID)
}

func (m *customRolesManager) GetRoleAssignmentsForMembership(ctx context.Context, membershipID string) ([]*customroles.CustomRoleAssignment, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repository.GetCustomRoleAssignmentsForMembership(ctx, membershipID)
}

func (m *customRolesManager) ListAvailablePermissions(_ context.Context) []string {
	allPerms := authorization.AllPermissionIDs()
	result := make([]string, 0, len(allPerms))
	for p := range allPerms {
		result = append(result, p)
	}
	sort.Strings(result)
	return result
}

package customroles

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/customroles"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/customroles/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
)

func (r *repository) CreateCustomRole(ctx context.Context, input *customroles.CustomRoleDatabaseCreationInput) (*customroles.CustomRole, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	tx, err := r.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}
	defer r.RollbackTransaction(ctx, tx)

	if err = r.generatedQuerier.CreateCustomRole(ctx, tx, &generated.CreateCustomRoleParams{
		ID:               input.ID,
		Name:             input.Name,
		Description:      input.Description,
		Scope:            input.Scope,
		BelongsToAccount: database.NullStringFromString(input.BelongsToAccount),
		CreatedBy:        input.CreatedBy,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating custom role")
	}

	for _, perm := range input.Permissions {
		if err = r.generatedQuerier.CreateCustomRolePermission(ctx, tx, &generated.CreateCustomRolePermissionParams{
			ID:            identifiers.New(),
			BelongsToRole: input.ID,
			Permission:    perm,
		}); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "creating custom role permission")
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return r.GetCustomRole(ctx, input.ID)
}

func (r *repository) GetCustomRole(ctx context.Context, roleID string) (*customroles.CustomRole, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if roleID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}

	row, err := r.generatedQuerier.GetCustomRole(ctx, r.readDB, roleID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting custom role")
	}

	perms, err := r.generatedQuerier.GetPermissionsForCustomRole(ctx, r.readDB, roleID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting permissions for custom role")
	}

	permissions := make([]string, len(perms))
	for i, p := range perms {
		permissions[i] = p.Permission
	}

	return convertRowToCustomRole(row, permissions), nil
}

func convertRowToCustomRole(row *generated.CustomRoles, permissions []string) *customroles.CustomRole {
	return &customroles.CustomRole{
		ID:               row.ID,
		Name:             row.Name,
		Description:      row.Description,
		Scope:            row.Scope,
		BelongsToAccount: row.BelongsToAccount.String,
		CreatedBy:        row.CreatedBy,
		CreatedAt:        row.CreatedAt,
		LastUpdatedAt:    database.TimePointerFromNullTime(row.LastUpdatedAt),
		ArchivedAt:       database.TimePointerFromNullTime(row.ArchivedAt),
		Permissions:      permissions,
	}
}

func (r *repository) ListServiceScopedCustomRoles(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[customroles.CustomRole], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	results, err := r.generatedQuerier.ListServiceScopedCustomRoles(ctx, r.readDB, &generated.ListServiceScopedCustomRolesParams{
		QueryLimit: database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "listing service-scoped custom roles")
	}

	var (
		data                      = []*customroles.CustomRole{}
		filteredCount, totalCount uint64
	)
	for _, row := range results {
		perms, permErr := r.generatedQuerier.GetPermissionsForCustomRole(ctx, r.readDB, row.ID)
		if permErr != nil {
			return nil, observability.PrepareAndLogError(permErr, logger, span, "getting permissions for custom role")
		}
		permissions := make([]string, len(perms))
		for i, p := range perms {
			permissions[i] = p.Permission
		}
		data = append(data, convertListServiceRowToCustomRole(row, permissions))
		filteredCount = uint64(row.FilteredCount)
		totalCount = uint64(row.TotalCount)
	}

	return filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *customroles.CustomRole) string { return t.ID },
		filter,
	), nil
}

func convertListServiceRowToCustomRole(row *generated.ListServiceScopedCustomRolesRow, permissions []string) *customroles.CustomRole {
	return &customroles.CustomRole{
		ID:               row.ID,
		Name:             row.Name,
		Description:      row.Description,
		Scope:            row.Scope,
		BelongsToAccount: row.BelongsToAccount.String,
		CreatedBy:        row.CreatedBy,
		CreatedAt:        row.CreatedAt,
		LastUpdatedAt:    database.TimePointerFromNullTime(row.LastUpdatedAt),
		ArchivedAt:       database.TimePointerFromNullTime(row.ArchivedAt),
		Permissions:      permissions,
	}
}

func (r *repository) ListAccountScopedCustomRoles(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[customroles.CustomRole], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if accountID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	results, err := r.generatedQuerier.ListAccountScopedCustomRoles(ctx, r.readDB, &generated.ListAccountScopedCustomRolesParams{
		BelongsToAccount: database.NullStringFromString(accountID),
		QueryLimit:       database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "listing account-scoped custom roles")
	}

	var (
		data                      = []*customroles.CustomRole{}
		filteredCount, totalCount uint64
	)
	for _, row := range results {
		perms, permErr := r.generatedQuerier.GetPermissionsForCustomRole(ctx, r.readDB, row.ID)
		if permErr != nil {
			return nil, observability.PrepareAndLogError(permErr, logger, span, "getting permissions for custom role")
		}
		permissions := make([]string, len(perms))
		for i, p := range perms {
			permissions[i] = p.Permission
		}
		data = append(data, convertListAccountRowToCustomRole(row, permissions))
		filteredCount = uint64(row.FilteredCount)
		totalCount = uint64(row.TotalCount)
	}

	return filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *customroles.CustomRole) string { return t.ID },
		filter,
	), nil
}

func convertListAccountRowToCustomRole(row *generated.ListAccountScopedCustomRolesRow, permissions []string) *customroles.CustomRole {
	return &customroles.CustomRole{
		ID:               row.ID,
		Name:             row.Name,
		Description:      row.Description,
		Scope:            row.Scope,
		BelongsToAccount: row.BelongsToAccount.String,
		CreatedBy:        row.CreatedBy,
		CreatedAt:        row.CreatedAt,
		LastUpdatedAt:    database.TimePointerFromNullTime(row.LastUpdatedAt),
		ArchivedAt:       database.TimePointerFromNullTime(row.ArchivedAt),
		Permissions:      permissions,
	}
}

func (r *repository) UpdateCustomRole(ctx context.Context, roleID string, input *customroles.CustomRoleUpdateRequestInput) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if roleID == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	existing, err := r.GetCustomRole(ctx, roleID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting existing custom role")
	}

	name := existing.Name
	if input.Name != nil {
		name = *input.Name
	}
	description := existing.Description
	if input.Description != nil {
		description = *input.Description
	}

	tx, err := r.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}
	defer r.RollbackTransaction(ctx, tx)

	if err = r.generatedQuerier.UpdateCustomRole(ctx, tx, &generated.UpdateCustomRoleParams{
		ID:          roleID,
		Name:        name,
		Description: description,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating custom role")
	}

	if input.Permissions != nil {
		// Delete existing permissions (raw SQL to bypass no-delete rule).
		if _, err = tx.ExecContext(ctx, "DELETE FROM custom_role_permissions WHERE belongs_to_role = $1", roleID); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "deleting existing permissions")
		}

		for _, perm := range input.Permissions {
			if err = r.generatedQuerier.CreateCustomRolePermission(ctx, tx, &generated.CreateCustomRolePermissionParams{
				ID:            identifiers.New(),
				BelongsToRole: roleID,
				Permission:    perm,
			}); err != nil {
				return observability.PrepareAndLogError(err, logger, span, "creating custom role permission")
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return nil
}

func (r *repository) ArchiveCustomRole(ctx context.Context, roleID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if roleID == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	return r.generatedQuerier.ArchiveCustomRole(ctx, r.writeDB, roleID)
}

func (r *repository) AssignCustomRoleToUser(ctx context.Context, id, roleID, userID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return r.generatedQuerier.AssignCustomRoleToUser(ctx, r.writeDB, &generated.AssignCustomRoleToUserParams{
		ID:           id,
		CustomRoleID: roleID,
		UserID:       database.NullStringFromString(userID),
	})
}

func (r *repository) AssignCustomRoleToMembership(ctx context.Context, id, roleID, membershipID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return r.generatedQuerier.AssignCustomRoleToMembership(ctx, r.writeDB, &generated.AssignCustomRoleToMembershipParams{
		ID:                  id,
		CustomRoleID:        roleID,
		AccountMembershipID: database.NullStringFromString(membershipID),
	})
}

func (r *repository) UnassignCustomRole(ctx context.Context, assignmentID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if assignmentID == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	return r.generatedQuerier.UnassignCustomRole(ctx, r.writeDB, assignmentID)
}

func (r *repository) GetCustomRoleAssignmentsForUser(ctx context.Context, userID string) ([]*customroles.CustomRoleAssignment, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	rows, err := r.generatedQuerier.GetCustomRoleAssignmentsForUser(ctx, r.readDB, database.NullStringFromString(userID))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting custom role assignments for user")
	}

	assignments := make([]*customroles.CustomRoleAssignment, len(rows))
	for i, row := range rows {
		assignments[i] = &customroles.CustomRoleAssignment{
			ID:                  row.ID,
			CustomRoleID:        row.CustomRoleID,
			CustomRoleName:      row.CustomRoleName,
			UserID:              row.UserID.String,
			AccountMembershipID: row.AccountMembershipID.String,
			CreatedAt:           row.CreatedAt,
		}
	}

	return assignments, nil
}

func (r *repository) GetCustomRoleAssignmentsForMembership(ctx context.Context, membershipID string) ([]*customroles.CustomRoleAssignment, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	rows, err := r.generatedQuerier.GetCustomRoleAssignmentsForMembership(ctx, r.readDB, database.NullStringFromString(membershipID))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting custom role assignments for membership")
	}

	assignments := make([]*customroles.CustomRoleAssignment, len(rows))
	for i, row := range rows {
		assignments[i] = &customroles.CustomRoleAssignment{
			ID:                  row.ID,
			CustomRoleID:        row.CustomRoleID,
			CustomRoleName:      row.CustomRoleName,
			UserID:              row.UserID.String,
			AccountMembershipID: row.AccountMembershipID.String,
			CreatedAt:           row.CreatedAt,
		}
	}

	return assignments, nil
}

func (r *repository) GetServiceScopedRoleIDsForUser(ctx context.Context, userID string) ([]string, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return r.generatedQuerier.GetServiceScopedRoleIDsForUser(ctx, r.readDB, database.NullStringFromString(userID))
}

func (r *repository) GetAccountScopedRoleIDsForMembership(ctx context.Context, membershipID string) ([]string, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return r.generatedQuerier.GetAccountScopedRoleIDsForMembership(ctx, r.readDB, database.NullStringFromString(membershipID))
}

func (r *repository) GetAllCustomRolePermissions(ctx context.Context) (map[string][]string, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	rows, err := r.generatedQuerier.GetAllCustomRolePermissions(ctx, r.readDB)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting all custom role permissions")
	}

	result := make(map[string][]string)
	for _, row := range rows {
		result[row.BelongsToRole] = append(result[row.BelongsToRole], row.Permission)
	}

	return result, nil
}

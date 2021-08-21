package postgres

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/Masterminds/squirrel"
)

var (
	_ querybuilding.HouseholdUserMembershipSQLQueryBuilder = (*Postgres)(nil)
)

const (
	householdMemberRolesSeparator = ","
)

// BuildGetDefaultHouseholdIDForUserQuery does .
func (b *Postgres) BuildGetDefaultHouseholdIDForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)

	return b.buildQuery(
		span,
		b.sqlBuilder.
			Select(fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.IDColumn)).
			From(querybuilding.HouseholdsTableName).
			Join(fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				querybuilding.HouseholdsUserMembershipTableName,
				querybuilding.HouseholdsUserMembershipTableName,
				querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn,
				querybuilding.HouseholdsTableName,
				querybuilding.IDColumn,
			)).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn):        userID,
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.HouseholdsUserMembershipTableDefaultUserHouseholdColumn): true,
			}),
	)
}

// BuildGetHouseholdMembershipsForUserQuery does .
func (b *Postgres) BuildGetHouseholdMembershipsForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)

	return b.buildQuery(
		span,
		b.sqlBuilder.
			Select(querybuilding.HouseholdsUserMembershipTableColumns...).
			Join(fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				querybuilding.HouseholdsTableName,
				querybuilding.HouseholdsTableName,
				querybuilding.IDColumn,
				querybuilding.HouseholdsUserMembershipTableName,
				querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn,
			)).
			From(querybuilding.HouseholdsUserMembershipTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.ArchivedOnColumn):                                 nil,
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn): userID,
			}),
	)
}

// BuildUserIsMemberOfHouseholdQuery builds a query that checks to see if the user is the member of a given household.
func (b *Postgres) BuildUserIsMemberOfHouseholdQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.
			Select(fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.HouseholdsUserMembershipTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn): householdID,
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn):      userID,
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.ArchivedOnColumn):                                      nil,
			}).
			Suffix(querybuilding.ExistenceSuffix),
	)
}

// BuildAddUserToHouseholdQuery builds a query that adds a user to an household.
func (b *Postgres) BuildAddUserToHouseholdQuery(ctx context.Context, input *types.AddUserToHouseholdInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, input.UserID)
	tracing.AttachHouseholdIDToSpan(span, input.HouseholdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.HouseholdsUserMembershipTableName).
			Columns(
				querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn,
				querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn,
				querybuilding.HouseholdsUserMembershipTableHouseholdRolesColumn,
			).
			Values(
				input.UserID,
				input.HouseholdID,
				strings.Join(input.HouseholdRoles, householdMemberRolesSeparator),
			),
	)
}

// BuildMarkHouseholdAsUserDefaultQuery builds a query that marks a user's household as their primary.
func (b *Postgres) BuildMarkHouseholdAsUserDefaultQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.HouseholdsUserMembershipTableName).
			Set(
				querybuilding.HouseholdsUserMembershipTableDefaultUserHouseholdColumn,
				squirrel.And{
					squirrel.Eq{querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn: userID},
					squirrel.Eq{querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn: householdID},
				},
			).
			Where(squirrel.Eq{
				querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn: userID,
				querybuilding.ArchivedOnColumn:                                 nil,
			}),
	)
}

// BuildModifyUserPermissionsQuery builds.
func (b *Postgres) BuildModifyUserPermissionsQuery(ctx context.Context, userID, householdID uint64, newRoles []string) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.HouseholdsUserMembershipTableName).
			Set(querybuilding.HouseholdsUserMembershipTableHouseholdRolesColumn, strings.Join(newRoles, householdMemberRolesSeparator)).
			Where(squirrel.Eq{
				querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn:      userID,
				querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn: householdID,
			}),
	)
}

// BuildTransferHouseholdOwnershipQuery does .
func (b *Postgres) BuildTransferHouseholdOwnershipQuery(ctx context.Context, currentOwnerID, newOwnerID, householdID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, newOwnerID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.HouseholdsTableName).
			Set(querybuilding.HouseholdsTableUserOwnershipColumn, newOwnerID).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                           householdID,
				querybuilding.HouseholdsTableUserOwnershipColumn: currentOwnerID,
				querybuilding.ArchivedOnColumn:                   nil,
			}),
	)
}

// BuildTransferHouseholdMembershipsQuery does .
func (b *Postgres) BuildTransferHouseholdMembershipsQuery(ctx context.Context, currentOwnerID, newOwnerID, householdID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, newOwnerID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.HouseholdsUserMembershipTableName).
			Set(querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn, newOwnerID).
			Where(squirrel.Eq{
				querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn: householdID,
				querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn:      currentOwnerID,
				querybuilding.ArchivedOnColumn:                                      nil,
			}),
	)
}

// BuildCreateMembershipForNewUserQuery builds a query that .
func (b *Postgres) BuildCreateMembershipForNewUserQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.HouseholdsUserMembershipTableName).
			Columns(
				querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn,
				querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn,
				querybuilding.HouseholdsUserMembershipTableDefaultUserHouseholdColumn,
				querybuilding.HouseholdsUserMembershipTableHouseholdRolesColumn,
			).
			Values(
				userID,
				householdID,
				true,
				strings.Join([]string{authorization.HouseholdAdminRole.String()}, householdMemberRolesSeparator),
			),
	)
}

// BuildRemoveUserFromHouseholdQuery builds a query that removes a user from an household.
func (b *Postgres) BuildRemoveUserFromHouseholdQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Delete(querybuilding.HouseholdsUserMembershipTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn): householdID,
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn):      userID,
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsUserMembershipTableName, querybuilding.ArchivedOnColumn):                                      nil,
			}),
	)
}

// BuildArchiveHouseholdMembershipsForUserQuery does .
func (b *Postgres) BuildArchiveHouseholdMembershipsForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.HouseholdsUserMembershipTableName).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.HouseholdsUserMembershipTableUserOwnershipColumn: userID,
				querybuilding.ArchivedOnColumn:                                 nil,
			}),
	)
}

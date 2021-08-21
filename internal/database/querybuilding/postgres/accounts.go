package postgres

import (
	"context"
	"fmt"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/Masterminds/squirrel"
)

var (
	_ querybuilding.HouseholdSQLQueryBuilder = (*Postgres)(nil)
)

// BuildGetHouseholdQuery constructs a SQL query for fetching an household with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetHouseholdQuery(ctx context.Context, householdID, userID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	columns := append(querybuilding.HouseholdsTableColumns, querybuilding.HouseholdsUserMembershipTableColumns...)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(columns...).
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
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.IDColumn):                           householdID,
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.HouseholdsTableUserOwnershipColumn): userID,
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.ArchivedOnColumn):                   nil,
			}),
	)
}

// BuildGetAllHouseholdsCountQuery returns a query that fetches the total number of households in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllHouseholdsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(span, b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.HouseholdsTableName)).
		From(querybuilding.HouseholdsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.ArchivedOnColumn): nil,
		}))
}

// BuildGetBatchOfHouseholdsQuery returns a query that fetches every household in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfHouseholdsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.HouseholdsTableColumns...).
			From(querybuilding.HouseholdsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetHouseholdsQuery builds a SQL query selecting households that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetHouseholdsQuery(ctx context.Context, userID uint64, forAdmin bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	var includeArchived bool
	if filter != nil {
		includeArchived = filter.IncludeArchived
	}

	columns := append(querybuilding.HouseholdsTableColumns, querybuilding.HouseholdsUserMembershipTableColumns...)
	filteredCountQuery, filteredCountQueryArgs := b.buildFilteredCountQuery(ctx, querybuilding.HouseholdsTableName, nil, nil, querybuilding.HouseholdsTableUserOwnershipColumn, userID, forAdmin, includeArchived, filter)
	totalCountQuery, totalCountQueryArgs := b.buildTotalCountQuery(ctx, querybuilding.HouseholdsTableName, nil, nil, querybuilding.HouseholdsTableUserOwnershipColumn, userID, forAdmin, includeArchived)

	builder := b.sqlBuilder.Select(append(
		columns,
		fmt.Sprintf("(%s) as total_count", totalCountQuery),
		fmt.Sprintf("(%s) as filtered_count", filteredCountQuery),
	)...).
		From(querybuilding.HouseholdsTableName).
		Join(fmt.Sprintf(
			"%s ON %s.%s = %s.%s",
			querybuilding.HouseholdsUserMembershipTableName,
			querybuilding.HouseholdsUserMembershipTableName,
			querybuilding.HouseholdsUserMembershipTableHouseholdOwnershipColumn,
			querybuilding.HouseholdsTableName,
			querybuilding.IDColumn,
		))

	if !forAdmin {
		builder = builder.Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.ArchivedOnColumn):                   nil,
			fmt.Sprintf("%s.%s", querybuilding.HouseholdsTableName, querybuilding.HouseholdsTableUserOwnershipColumn): userID,
		})
	}

	builder = builder.GroupBy(fmt.Sprintf(
		"%s.%s, %s.%s",
		querybuilding.HouseholdsTableName,
		querybuilding.IDColumn,
		querybuilding.HouseholdsUserMembershipTableName,
		querybuilding.IDColumn,
	))

	if filter != nil {
		builder = querybuilding.ApplyFilterToQueryBuilder(filter, querybuilding.HouseholdsTableName, builder)
	}

	query, selectArgs := b.buildQuery(span, builder)

	return query, append(append(filteredCountQueryArgs, totalCountQueryArgs...), selectArgs...)
}

// BuildHouseholdCreationQuery takes an household and returns a creation query for that household and the relevant arguments.
func (b *Postgres) BuildHouseholdCreationQuery(ctx context.Context, input *types.HouseholdCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.HouseholdsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.HouseholdsTableNameColumn,
				querybuilding.HouseholdsTableBillingStatusColumn,
				querybuilding.HouseholdsTableContactEmailColumn,
				querybuilding.HouseholdsTableContactPhoneColumn,
				querybuilding.HouseholdsTableUserOwnershipColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Name,
				types.UnpaidHouseholdBillingStatus,
				input.ContactEmail,
				input.ContactPhone,
				input.BelongsToUser,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateHouseholdQuery takes an household and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateHouseholdQuery(ctx context.Context, input *types.Household) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachHouseholdIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.HouseholdsTableName).
			Set(querybuilding.HouseholdsTableNameColumn, input.Name).
			Set(querybuilding.HouseholdsTableContactEmailColumn, input.ContactEmail).
			Set(querybuilding.HouseholdsTableContactPhoneColumn, input.ContactPhone).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                           input.ID,
				querybuilding.ArchivedOnColumn:                   nil,
				querybuilding.HouseholdsTableUserOwnershipColumn: input.BelongsToUser,
			}),
	)
}

// BuildArchiveHouseholdQuery returns a SQL query which marks a given household belonging to a given user as archived.
func (b *Postgres) BuildArchiveHouseholdQuery(ctx context.Context, householdID, userID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.HouseholdsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                           householdID,
				querybuilding.ArchivedOnColumn:                   nil,
				querybuilding.HouseholdsTableUserOwnershipColumn: userID,
			}),
	)
}

// BuildGetAuditLogEntriesForHouseholdQuery constructs a SQL query for fetching audit log entries belong to an household with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForHouseholdQuery(ctx context.Context, householdID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachHouseholdIDToSpan(span, householdID)

	householdIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.HouseholdAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{householdIDKey: householdID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}

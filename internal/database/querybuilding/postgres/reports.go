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

var _ querybuilding.ReportSQLQueryBuilder = (*Postgres)(nil)

// BuildReportExistsQuery constructs a SQL query for checking if a report with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildReportExistsQuery(ctx context.Context, reportID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachReportIDToSpan(span, reportID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.ReportsTableName).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.IDColumn):         reportID,
				fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetReportQuery constructs a SQL query for fetching a report with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachReportIDToSpan(span, reportID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ReportsTableColumns...).
			From(querybuilding.ReportsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.IDColumn):         reportID,
				fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetAllReportsCountQuery returns a query that fetches the total number of reports in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllReportsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.ReportsTableName)).
			From(querybuilding.ReportsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfReportsQuery returns a query that fetches every report in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfReportsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ReportsTableColumns...).
			From(querybuilding.ReportsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetReportsQuery builds a SQL query selecting reports that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetReportsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.ReportsTableName,
		nil,
		where,
		"",
		querybuilding.ReportsTableColumns,
		0,
		includeArchived,
		filter,
	)
}

// BuildGetReportsWithIDsQuery builds a SQL query selecting reports that belong to a given household,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetReportsWithIDsQuery(ctx context.Context, householdID uint64, limit uint8, ids []uint64, restrictToHousehold bool) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.IDColumn):         ids,
		fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	if restrictToHousehold {
		where[fmt.Sprintf("%s.%s", querybuilding.ReportsTableName, querybuilding.ReportsTableHouseholdOwnershipColumn)] = householdID
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.ReportsTableColumns...).
		From(querybuilding.ReportsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ReportsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.ReportsTableName).
			Where(where),
	)
}

// BuildCreateReportQuery takes a report and returns a creation query for that report and the relevant arguments.
func (b *Postgres) BuildCreateReportQuery(ctx context.Context, input *types.ReportCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.ReportsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.ReportsTableReportTypeColumn,
				querybuilding.ReportsTableConcernColumn,
				querybuilding.ReportsTableHouseholdOwnershipColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.ReportType,
				input.Concern,
				input.BelongsToHousehold,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateReportQuery takes a report and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateReportQuery(ctx context.Context, input *types.Report) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachReportIDToSpan(span, input.ID)
	tracing.AttachHouseholdIDToSpan(span, input.BelongsToHousehold)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ReportsTableName).
			Set(querybuilding.ReportsTableReportTypeColumn, input.ReportType).
			Set(querybuilding.ReportsTableConcernColumn, input.Concern).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                             input.ID,
				querybuilding.ArchivedOnColumn:                     nil,
				querybuilding.ReportsTableHouseholdOwnershipColumn: input.BelongsToHousehold,
			}),
	)
}

// BuildArchiveReportQuery returns a SQL query which marks a given report belonging to a given household as archived.
func (b *Postgres) BuildArchiveReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachReportIDToSpan(span, reportID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ReportsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         reportID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildGetAuditLogEntriesForReportQuery constructs a SQL query for fetching audit log entries relating to a report with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachReportIDToSpan(span, reportID)

	reportIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.ReportAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{reportIDKey: reportID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}

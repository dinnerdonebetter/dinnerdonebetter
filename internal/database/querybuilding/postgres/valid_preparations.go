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

var _ querybuilding.ValidPreparationSQLQueryBuilder = (*Postgres)(nil)

// BuildValidPreparationExistsQuery constructs a SQL query for checking if a valid preparation with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildValidPreparationExistsQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.ValidPreparationsTableName).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.IDColumn):         validPreparationID,
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetValidPreparationQuery constructs a SQL query for fetching a valid preparation with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidPreparationsTableColumns...).
			From(querybuilding.ValidPreparationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.IDColumn):         validPreparationID,
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetAllValidPreparationsCountQuery returns a query that fetches the total number of valid preparations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllValidPreparationsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.ValidPreparationsTableName)).
			From(querybuilding.ValidPreparationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfValidPreparationsQuery returns a query that fetches every valid preparation in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfValidPreparationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidPreparationsTableColumns...).
			From(querybuilding.ValidPreparationsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetValidPreparationsQuery builds a SQL query selecting valid preparations that adhere to a given QueryFilter and belong to a given account,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetValidPreparationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.ValidPreparationsTableName,
		nil,
		where,
		"",
		querybuilding.ValidPreparationsTableColumns,
		0,
		includeArchived,
		filter,
	)
}

// BuildGetValidPreparationsWithIDsQuery builds a SQL query selecting valid preparations that belong to a given account,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetValidPreparationsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.IDColumn):         ids,
		fmt.Sprintf("%s.%s", querybuilding.ValidPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.ValidPreparationsTableColumns...).
		From(querybuilding.ValidPreparationsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidPreparationsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.ValidPreparationsTableName).
			Where(where),
	)
}

// BuildCreateValidPreparationQuery takes a valid preparation and returns a creation query for that valid preparation and the relevant arguments.
func (b *Postgres) BuildCreateValidPreparationQuery(ctx context.Context, input *types.ValidPreparationCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.ValidPreparationsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.ValidPreparationsTableNameColumn,
				querybuilding.ValidPreparationsTableDescriptionColumn,
				querybuilding.ValidPreparationsTableIconPathColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Name,
				input.Description,
				input.IconPath,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateValidPreparationQuery takes a valid preparation and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateValidPreparationQuery(ctx context.Context, input *types.ValidPreparation) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidPreparationsTableName).
			Set(querybuilding.ValidPreparationsTableNameColumn, input.Name).
			Set(querybuilding.ValidPreparationsTableDescriptionColumn, input.Description).
			Set(querybuilding.ValidPreparationsTableIconPathColumn, input.IconPath).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         input.ID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildArchiveValidPreparationQuery returns a SQL query which marks a given valid preparation belonging to a given account as archived.
func (b *Postgres) BuildArchiveValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidPreparationsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         validPreparationID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildGetAuditLogEntriesForValidPreparationQuery constructs a SQL query for fetching audit log entries relating to a valid preparation with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	validPreparationIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.ValidPreparationAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{validPreparationIDKey: validPreparationID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}

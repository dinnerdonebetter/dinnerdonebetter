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

var _ querybuilding.ValidPreparationInstrumentSQLQueryBuilder = (*Postgres)(nil)

// BuildValidPreparationInstrumentExistsQuery constructs a SQL query for checking if a valid preparation instrument with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildValidPreparationInstrumentExistsQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.ValidPreparationInstrumentsTableName).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.IDColumn):         validPreparationInstrumentID,
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetValidPreparationInstrumentQuery constructs a SQL query for fetching a valid preparation instrument with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidPreparationInstrumentsTableColumns...).
			From(querybuilding.ValidPreparationInstrumentsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.IDColumn):         validPreparationInstrumentID,
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetAllValidPreparationInstrumentsCountQuery returns a query that fetches the total number of valid preparation instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllValidPreparationInstrumentsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.ValidPreparationInstrumentsTableName)).
			From(querybuilding.ValidPreparationInstrumentsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfValidPreparationInstrumentsQuery returns a query that fetches every valid preparation instrument in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfValidPreparationInstrumentsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidPreparationInstrumentsTableColumns...).
			From(querybuilding.ValidPreparationInstrumentsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetValidPreparationInstrumentsQuery builds a SQL query selecting valid preparation instruments that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetValidPreparationInstrumentsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.ValidPreparationInstrumentsTableName,
		nil,
		where,
		"",
		querybuilding.ValidPreparationInstrumentsTableColumns,
		0,
		includeArchived,
		filter,
	)
}

// BuildGetValidPreparationInstrumentsWithIDsQuery builds a SQL query selecting valid preparation instruments that belong to a given household,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetValidPreparationInstrumentsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.IDColumn):         ids,
		fmt.Sprintf("%s.%s", querybuilding.ValidPreparationInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.ValidPreparationInstrumentsTableColumns...).
		From(querybuilding.ValidPreparationInstrumentsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidPreparationInstrumentsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.ValidPreparationInstrumentsTableName).
			Where(where),
	)
}

// BuildCreateValidPreparationInstrumentQuery takes a valid preparation instrument and returns a creation query for that valid preparation instrument and the relevant arguments.
func (b *Postgres) BuildCreateValidPreparationInstrumentQuery(ctx context.Context, input *types.ValidPreparationInstrumentCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.ValidPreparationInstrumentsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.ValidPreparationInstrumentsTableInstrumentIDColumn,
				querybuilding.ValidPreparationInstrumentsTablePreparationIDColumn,
				querybuilding.ValidPreparationInstrumentsTableNotesColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.InstrumentID,
				input.PreparationID,
				input.Notes,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateValidPreparationInstrumentQuery takes a valid preparation instrument and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateValidPreparationInstrumentQuery(ctx context.Context, input *types.ValidPreparationInstrument) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationInstrumentIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidPreparationInstrumentsTableName).
			Set(querybuilding.ValidPreparationInstrumentsTableInstrumentIDColumn, input.InstrumentID).
			Set(querybuilding.ValidPreparationInstrumentsTablePreparationIDColumn, input.PreparationID).
			Set(querybuilding.ValidPreparationInstrumentsTableNotesColumn, input.Notes).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         input.ID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildArchiveValidPreparationInstrumentQuery returns a SQL query which marks a given valid preparation instrument belonging to a given household as archived.
func (b *Postgres) BuildArchiveValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidPreparationInstrumentsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         validPreparationInstrumentID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildGetAuditLogEntriesForValidPreparationInstrumentQuery constructs a SQL query for fetching audit log entries relating to a valid preparation instrument with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	validPreparationInstrumentIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.ValidPreparationInstrumentAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{validPreparationInstrumentIDKey: validPreparationInstrumentID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}

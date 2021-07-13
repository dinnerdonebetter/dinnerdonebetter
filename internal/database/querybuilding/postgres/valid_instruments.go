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

var _ querybuilding.ValidInstrumentSQLQueryBuilder = (*Postgres)(nil)

// BuildValidInstrumentExistsQuery constructs a SQL query for checking if a valid instrument with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildValidInstrumentExistsQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.ValidInstrumentsTableName).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.IDColumn):         validInstrumentID,
				fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetValidInstrumentQuery constructs a SQL query for fetching a valid instrument with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidInstrumentsTableColumns...).
			From(querybuilding.ValidInstrumentsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.IDColumn):         validInstrumentID,
				fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetAllValidInstrumentsCountQuery returns a query that fetches the total number of valid instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllValidInstrumentsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.ValidInstrumentsTableName)).
			From(querybuilding.ValidInstrumentsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfValidInstrumentsQuery returns a query that fetches every valid instrument in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfValidInstrumentsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidInstrumentsTableColumns...).
			From(querybuilding.ValidInstrumentsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetValidInstrumentsQuery builds a SQL query selecting valid instruments that adhere to a given QueryFilter and belong to a given account,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetValidInstrumentsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.ValidInstrumentsTableName,
		nil,
		where,
		"",
		querybuilding.ValidInstrumentsTableColumns,
		0,
		includeArchived,
		filter,
	)
}

// BuildGetValidInstrumentsWithIDsQuery builds a SQL query selecting valid instruments that belong to a given account,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetValidInstrumentsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.IDColumn):         ids,
		fmt.Sprintf("%s.%s", querybuilding.ValidInstrumentsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.ValidInstrumentsTableColumns...).
		From(querybuilding.ValidInstrumentsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidInstrumentsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.ValidInstrumentsTableName).
			Where(where),
	)
}

// BuildCreateValidInstrumentQuery takes a valid instrument and returns a creation query for that valid instrument and the relevant arguments.
func (b *Postgres) BuildCreateValidInstrumentQuery(ctx context.Context, input *types.ValidInstrumentCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.ValidInstrumentsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.ValidInstrumentsTableNameColumn,
				querybuilding.ValidInstrumentsTableVariantColumn,
				querybuilding.ValidInstrumentsTableDescriptionColumn,
				querybuilding.ValidInstrumentsTableIconPathColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Name,
				input.Variant,
				input.Description,
				input.IconPath,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateValidInstrumentQuery takes a valid instrument and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateValidInstrumentQuery(ctx context.Context, input *types.ValidInstrument) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidInstrumentsTableName).
			Set(querybuilding.ValidInstrumentsTableNameColumn, input.Name).
			Set(querybuilding.ValidInstrumentsTableVariantColumn, input.Variant).
			Set(querybuilding.ValidInstrumentsTableDescriptionColumn, input.Description).
			Set(querybuilding.ValidInstrumentsTableIconPathColumn, input.IconPath).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         input.ID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildArchiveValidInstrumentQuery returns a SQL query which marks a given valid instrument belonging to a given account as archived.
func (b *Postgres) BuildArchiveValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidInstrumentsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         validInstrumentID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildGetAuditLogEntriesForValidInstrumentQuery constructs a SQL query for fetching audit log entries relating to a valid instrument with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	validInstrumentIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.ValidInstrumentAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{validInstrumentIDKey: validInstrumentID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}

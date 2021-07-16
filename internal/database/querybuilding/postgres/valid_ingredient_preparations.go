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

var _ querybuilding.ValidIngredientPreparationSQLQueryBuilder = (*Postgres)(nil)

// BuildValidIngredientPreparationExistsQuery constructs a SQL query for checking if a valid ingredient preparation with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildValidIngredientPreparationExistsQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.ValidIngredientPreparationsTableName).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.IDColumn):         validIngredientPreparationID,
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetValidIngredientPreparationQuery constructs a SQL query for fetching a valid ingredient preparation with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidIngredientPreparationsTableColumns...).
			From(querybuilding.ValidIngredientPreparationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.IDColumn):         validIngredientPreparationID,
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetAllValidIngredientPreparationsCountQuery returns a query that fetches the total number of valid ingredient preparations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllValidIngredientPreparationsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.ValidIngredientPreparationsTableName)).
			From(querybuilding.ValidIngredientPreparationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfValidIngredientPreparationsQuery returns a query that fetches every valid ingredient preparation in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfValidIngredientPreparationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidIngredientPreparationsTableColumns...).
			From(querybuilding.ValidIngredientPreparationsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetValidIngredientPreparationsQuery builds a SQL query selecting valid ingredient preparations that adhere to a given QueryFilter and belong to a given account,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetValidIngredientPreparationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.ValidIngredientPreparationsTableName,
		nil,
		where,
		"",
		querybuilding.ValidIngredientPreparationsTableColumns,
		0,
		includeArchived,
		filter,
	)
}

// BuildGetValidIngredientPreparationsWithIDsQuery builds a SQL query selecting valid ingredient preparations that belong to a given account,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetValidIngredientPreparationsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.IDColumn):         ids,
		fmt.Sprintf("%s.%s", querybuilding.ValidIngredientPreparationsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.ValidIngredientPreparationsTableColumns...).
		From(querybuilding.ValidIngredientPreparationsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidIngredientPreparationsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.ValidIngredientPreparationsTableName).
			Where(where),
	)
}

// BuildCreateValidIngredientPreparationQuery takes a valid ingredient preparation and returns a creation query for that valid ingredient preparation and the relevant arguments.
func (b *Postgres) BuildCreateValidIngredientPreparationQuery(ctx context.Context, input *types.ValidIngredientPreparationCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.ValidIngredientPreparationsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.ValidIngredientPreparationsTableNotesColumn,
				querybuilding.ValidIngredientPreparationsTableValidIngredientIDColumn,
				querybuilding.ValidIngredientPreparationsTableValidPreparationIDColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Notes,
				input.ValidIngredientID,
				input.ValidPreparationID,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateValidIngredientPreparationQuery takes a valid ingredient preparation and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateValidIngredientPreparationQuery(ctx context.Context, input *types.ValidIngredientPreparation) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidIngredientPreparationsTableName).
			Set(querybuilding.ValidIngredientPreparationsTableNotesColumn, input.Notes).
			Set(querybuilding.ValidIngredientPreparationsTableValidIngredientIDColumn, input.ValidIngredientID).
			Set(querybuilding.ValidIngredientPreparationsTableValidPreparationIDColumn, input.ValidPreparationID).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         input.ID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildArchiveValidIngredientPreparationQuery returns a SQL query which marks a given valid ingredient preparation belonging to a given account as archived.
func (b *Postgres) BuildArchiveValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidIngredientPreparationsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         validIngredientPreparationID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildGetAuditLogEntriesForValidIngredientPreparationQuery constructs a SQL query for fetching audit log entries relating to a valid ingredient preparation with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	validIngredientPreparationIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.ValidIngredientPreparationAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{validIngredientPreparationIDKey: validIngredientPreparationID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	defaultLimit             = types.DefaultLimit
	commaSeparator           = ","
	columnCountQueryTemplate = `COUNT(%s.id)`
	userOwnershipColumn      = "belongs_to_user"
	householdOwnershipColumn = "belongs_to_household"
)

func wrapQueryForILIKE(s string) string {
	return fmt.Sprintf("%%%s%%", s)
}

func joinIDs(ids []string) string {
	return strings.Join(ids, commaSeparator)
}

// logQueryBuildingError logs errs that may occur during query construction. Such errors should be few and far between,
// as the generally only occur with type discrepancies or other misuses of SQL. An alert should be set up for any log
// entries with the given name, and those alerts should be investigated quickly.
func (q *SQLQuerier) logQueryBuildingError(span tracing.Span, err error) {
	if err != nil {
		logger := q.logger.WithValue(keys.QueryErrorKey, true)
		observability.AcknowledgeError(err, logger, span, "building query")
	}
}

// BuildQuery builds a given query, handles whatever errs and returns just the query and args.
func (q *SQLQuerier) buildQuery(span tracing.Span, builder squirrel.Sqlizer) (query string, args []interface{}) {
	query, args, err := builder.ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

func (q *SQLQuerier) buildTotalCountQuery(
	ctx context.Context,
	tableName string,
	joins []string,
	where squirrel.Eq,
	ownershipColumn,
	userID string,
	forAdmin,
	includeArchived bool,
) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if where == nil {
		where = squirrel.Eq{}
	}

	totalCountQueryBuilder := q.sqlBuilder.
		PlaceholderFormat(squirrel.Question).
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	for _, join := range joins {
		totalCountQueryBuilder = totalCountQueryBuilder.Join(join)
	}

	if !forAdmin {
		if userID != "" && ownershipColumn != "" {
			where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = userID
		}

		where[fmt.Sprintf("%s.archived_on", tableName)] = nil
	} else if !includeArchived {
		where[fmt.Sprintf("%s.archived_on", tableName)] = nil
	}

	if len(where) > 0 {
		totalCountQueryBuilder = totalCountQueryBuilder.Where(where)
	}

	return q.buildQuery(span, totalCountQueryBuilder)
}

func (q *SQLQuerier) buildTotalCountQueryWithILike(
	ctx context.Context,
	tableName string,
	joins []string,
	where squirrel.ILike,
	ownershipColumn,
	userID string,
	forAdmin,
	includeArchived bool,
) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	totalCountQueryBuilder := q.sqlBuilder.
		PlaceholderFormat(squirrel.Question).
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	for _, join := range joins {
		totalCountQueryBuilder = totalCountQueryBuilder.Join(join)
	}

	equalsWhere := squirrel.Eq{}
	if !forAdmin {
		if userID != "" && ownershipColumn != "" {
			equalsWhere[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = userID
		}

		equalsWhere[fmt.Sprintf("%s.archived_on", tableName)] = nil
	} else if !includeArchived {
		equalsWhere[fmt.Sprintf("%s.archived_on", tableName)] = nil
	}

	if where != nil {
		totalCountQueryBuilder = totalCountQueryBuilder.Where(where)
	}

	if len(equalsWhere) > 0 {
		totalCountQueryBuilder = totalCountQueryBuilder.Where(equalsWhere)
	}

	return q.buildQuery(span, totalCountQueryBuilder)
}

func (q *SQLQuerier) buildFilteredCountQuery(
	ctx context.Context,
	tableName string,
	joins []string,
	where squirrel.Eq,
	ownershipColumn,
	userID string,
	forAdmin,
	includeArchived bool,
	filter *types.QueryFilter,
) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)
	}

	if where == nil {
		where = squirrel.Eq{}
	}

	filteredCountQueryBuilder := q.sqlBuilder.
		PlaceholderFormat(squirrel.Question).
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	for _, join := range joins {
		filteredCountQueryBuilder = filteredCountQueryBuilder.Join(join)
	}

	if !forAdmin {
		if userID != "" && ownershipColumn != "" {
			where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = userID
		}

		where[fmt.Sprintf("%s.archived_on", tableName)] = nil
	} else if !includeArchived {
		where[fmt.Sprintf("%s.archived_on", tableName)] = nil
	}

	if len(where) > 0 {
		filteredCountQueryBuilder = filteredCountQueryBuilder.Where(where)
	}

	if filter != nil {
		filteredCountQueryBuilder = applyFilterToSubCountQueryBuilder(filter, tableName, filteredCountQueryBuilder)
	}

	return q.buildQuery(span, filteredCountQueryBuilder)
}

func (q *SQLQuerier) buildFilteredCountQueryWithILike(
	ctx context.Context,
	tableName string,
	joins []string,
	where squirrel.ILike,
	ownershipColumn,
	userID string,
	forAdmin,
	includeArchived bool,
	filter *types.QueryFilter,
) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)
	}

	filteredCountQueryBuilder := q.sqlBuilder.
		PlaceholderFormat(squirrel.Question).
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName).
		Where(where)

	for _, join := range joins {
		filteredCountQueryBuilder = filteredCountQueryBuilder.Join(join)
	}

	equalsWhere := squirrel.Eq{}
	if !forAdmin {
		if userID != "" && ownershipColumn != "" {
			equalsWhere[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = userID
		}

		equalsWhere[fmt.Sprintf("%s.archived_on", tableName)] = nil
	} else if !includeArchived {
		equalsWhere[fmt.Sprintf("%s.archived_on", tableName)] = nil
	}

	if len(equalsWhere) > 0 {
		filteredCountQueryBuilder = filteredCountQueryBuilder.Where(equalsWhere)
	}

	if filter != nil {
		filteredCountQueryBuilder = applyFilterToSubCountQueryBuilder(filter, tableName, filteredCountQueryBuilder)
	}

	return q.buildQuery(span, filteredCountQueryBuilder)
}

// BuildListQuery builds a SQL query selecting rows that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (q *SQLQuerier) buildListQuery(
	ctx context.Context,
	tableName string,
	joins,
	groupBys []string,
	where squirrel.Eq,
	ownershipColumn string,
	columns []string,
	ownerID string,
	forAdmin bool,
	filter *types.QueryFilter,
	includeCountQueries bool,
) (query string, args []interface{}) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)
	}

	var includeArchived bool
	if filter != nil {
		if filter.IncludeArchived != nil {
			includeArchived = *filter.IncludeArchived
		}
	}

	filteredCountQuery, filteredCountQueryArgs := q.buildFilteredCountQuery(ctx, tableName, joins, where, ownershipColumn, ownerID, forAdmin, includeArchived, filter)
	totalCountQuery, totalCountQueryArgs := q.buildTotalCountQuery(ctx, tableName, joins, where, ownershipColumn, ownerID, forAdmin, includeArchived)

	if includeCountQueries {
		columns = append(
			columns,
			fmt.Sprintf("(%s) as filtered_count", filteredCountQuery),
			fmt.Sprintf("(%s) as total_count", totalCountQuery),
		)
	}

	builder := q.sqlBuilder.Select(columns...).From(tableName)
	for _, join := range joins {
		builder = builder.Join(join)
	}

	if !forAdmin {
		if where == nil {
			where = squirrel.Eq{}
		}
		where[fmt.Sprintf("%s.archived_on", tableName)] = nil

		if ownershipColumn != "" && ownerID != "" {
			where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = ownerID
		}

		builder = builder.Where(where)
	}

	actualGroupBys := []string{fmt.Sprintf("%s.%s", tableName, "id")}
	if groupBys != nil {
		actualGroupBys = append(actualGroupBys, groupBys...)
	}

	builder = builder.GroupBy(actualGroupBys...)
	builder = builder.OrderBy(fmt.Sprintf("%s.%s", tableName, "id"))

	if filter != nil {
		builder = applyFilterToQueryBuilder(filter, tableName, builder)
	}

	query, args = q.buildQuery(span, builder)
	if includeCountQueries {
		args = append(append(filteredCountQueryArgs, totalCountQueryArgs...), args...)
	}

	return query, args
}

// BuildListQueryWithILike builds a SQL query selecting rows that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (q *SQLQuerier) buildListQueryWithILike(
	ctx context.Context,
	tableName string,
	joins,
	groupBys []string,
	where squirrel.ILike,
	ownershipColumn string,
	columns []string,
	ownerID string,
	forAdmin bool,
	filter *types.QueryFilter,
) (query string, args []interface{}) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)
	}

	var includeArchived bool
	if filter != nil {
		if filter.IncludeArchived != nil {
			includeArchived = *filter.IncludeArchived
		}
	}

	filteredCountQuery, filteredCountQueryArgs := q.buildFilteredCountQueryWithILike(ctx, tableName, joins, where, ownershipColumn, ownerID, forAdmin, includeArchived, filter)
	totalCountQuery, totalCountQueryArgs := q.buildTotalCountQueryWithILike(ctx, tableName, joins, where, ownershipColumn, ownerID, forAdmin, includeArchived)

	builder := q.sqlBuilder.
		Select(append(
			columns,
			fmt.Sprintf("(%s) as filtered_count", filteredCountQuery),
			fmt.Sprintf("(%s) as total_count", totalCountQuery),
		)...).
		From(tableName).
		Where(where)

	for _, join := range joins {
		builder = builder.Join(join)
	}

	equalsWhere := squirrel.Eq{}
	if !forAdmin {
		equalsWhere[fmt.Sprintf("%s.archived_on", tableName)] = nil

		if ownershipColumn != "" && ownerID != "" {
			equalsWhere[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = ownerID
		}
	}
	builder = builder.Where(equalsWhere)

	actualGroupBys := []string{fmt.Sprintf("%s.%s", tableName, "id")}
	if groupBys != nil {
		actualGroupBys = append(actualGroupBys, groupBys...)
	}

	builder = builder.GroupBy(actualGroupBys...)
	builder = builder.OrderBy(fmt.Sprintf("%s.%s", tableName, "id"))

	if filter != nil {
		builder = applyFilterToQueryBuilder(filter, tableName, builder)
	}

	query, selectArgs := q.buildQuery(span, builder)

	return query, append(append(filteredCountQueryArgs, totalCountQueryArgs...), selectArgs...)
}

func nullInt64ForUint64Field(x *uint64) sql.NullInt64 {
	z := sql.NullInt64{}

	if x != nil {
		z.Int64 = int64(*x)
		z.Valid = true
	}

	return z
}

func nullInt32ForUint8Field(x *uint8) sql.NullInt32 {
	z := sql.NullInt32{}

	if x != nil {
		z.Int32 = int32(*x)
		z.Valid = true
	}

	return z
}

package postgres

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/Masterminds/squirrel"
)

const (
	columnCountQueryTemplate = `COUNT(%s.id)`
	householdOwnershipColumn = "belongs_to_household"
)

// logQueryBuildingError logs errs that may occur during query construction. Such errors should be few and far between,
// as the generally only occur with type discrepancies or other misuses of SQL. An alert should be set up for any log
// entries with the given name, and those alerts should be investigated quickly.
func (q *Querier) logQueryBuildingError(span tracing.Span, err error) {
	if err != nil {
		logger := q.logger.WithValue(keys.QueryErrorKey, true)
		observability.AcknowledgeError(err, logger, span, "building query")
	}
}

// BuildQuery builds a given query, handles whatever errs and returns just the query and args.
func (q *Querier) buildQuery(span tracing.Span, builder squirrel.Sqlizer) (query string, args []any) {
	query, args, err := builder.ToSql()
	q.logQueryBuildingError(span, err)
	return query, args
}

func (q *Querier) buildTotalCountQuery(ctx context.Context, tableName string, joins []string, ownershipColumn, userID string, includeArchived bool) (query string, args []any) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{}

	totalCountQueryBuilder := q.sqlBuilder.
		PlaceholderFormat(squirrel.Question).
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	for _, join := range joins {
		totalCountQueryBuilder = totalCountQueryBuilder.Join(join)
	}

	if userID != "" && ownershipColumn != "" {
		where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = userID
	}

	where[fmt.Sprintf("%s.archived_at", tableName)] = nil

	if !includeArchived {
		where[fmt.Sprintf("%s.archived_at", tableName)] = nil
	}

	if len(where) > 0 {
		totalCountQueryBuilder = totalCountQueryBuilder.Where(where)
	}

	return q.buildQuery(span, totalCountQueryBuilder)
}

func (q *Querier) buildFilteredCountQuery(ctx context.Context, tableName string, joins []string, ownershipColumn, userID string, includeArchived bool, filter *types.QueryFilter) (query string, args []any) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)
	}

	where := squirrel.Eq{}

	filteredCountQueryBuilder := q.sqlBuilder.
		PlaceholderFormat(squirrel.Question).
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	for _, join := range joins {
		filteredCountQueryBuilder = filteredCountQueryBuilder.Join(join)
	}

	if userID != "" && ownershipColumn != "" {
		where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = userID
	}

	where[fmt.Sprintf("%s.archived_at", tableName)] = nil

	if !includeArchived {
		where[fmt.Sprintf("%s.archived_at", tableName)] = nil
	}

	if len(where) > 0 {
		filteredCountQueryBuilder = filteredCountQueryBuilder.Where(where)
	}

	if filter != nil {
		filteredCountQueryBuilder = applyFilterToSubCountQueryBuilder(filter, tableName, filteredCountQueryBuilder)
	}

	return q.buildQuery(span, filteredCountQueryBuilder)
}

// BuildListQuery builds a SQL query selecting rows that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor. TODO: Deprecate this.
func (q *Querier) buildListQuery(ctx context.Context, tableName string, joins, groupBys []string, ownershipColumn string, columns []string, ownerID string, filter *types.QueryFilter) (query string, args []any) {
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

	filteredCountQuery, filteredCountQueryArgs := q.buildFilteredCountQuery(ctx, tableName, joins, ownershipColumn, ownerID, includeArchived, filter)
	totalCountQuery, totalCountQueryArgs := q.buildTotalCountQuery(ctx, tableName, joins, ownershipColumn, ownerID, includeArchived)

	columns = append(
		columns,
		fmt.Sprintf("(%s) as filtered_count", filteredCountQuery),
		fmt.Sprintf("(%s) as total_count", totalCountQuery),
	)

	builder := q.sqlBuilder.Select(columns...).From(tableName)
	for _, join := range joins {
		builder = builder.LeftJoin(join)
	}

	where := squirrel.Eq{}
	where[fmt.Sprintf("%s.archived_at", tableName)] = nil

	if ownershipColumn != "" && ownerID != "" {
		where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = ownerID
	}

	builder = builder.Where(where)

	actualGroupBys := []string{fmt.Sprintf("%s.%s", tableName, "id")}
	for _, groupBy := range groupBys {
		if groupBy != actualGroupBys[0] {
			actualGroupBys = append(actualGroupBys, groupBy)
		}
	}

	builder = builder.GroupBy(actualGroupBys...).OrderBy(fmt.Sprintf("%s.%s", tableName, "id"))

	if filter != nil {
		builder = applyFilterToQueryBuilder(filter, tableName, builder)
	}

	query, args = q.buildQuery(span, builder)
	args = append(append(filteredCountQueryArgs, totalCountQueryArgs...), args...)

	return query, args
}

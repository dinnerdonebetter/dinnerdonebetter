package queries

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/Masterminds/squirrel"
)

type QuerySpec struct {
	where         any
	queryFilter   *types.QueryFilter
	table         string
	orderBy       string
	columns       []string
	joins         []string
	groupBys      []string
	limit         uint64
	offset        uint64
	includeCounts bool
}

func (s *QuerySpec) Table(table string) {
	s.table = table
}

func (s *QuerySpec) WithJoins(joins []string) {
	s.joins = joins
}

func (s *QuerySpec) WithWhere(where any) {
	s.where = where
}

func (s *QuerySpec) WithGroupBys(groupBys []string) {
	s.groupBys = groupBys
}

func (s *QuerySpec) WithOrderBy(orderBy string) {
	s.orderBy = orderBy
}

func (s *QuerySpec) WithLimit(limit uint64) {
	s.limit = limit
}

func (s *QuerySpec) WithOffset(offset uint64) {
	s.offset = offset
}

func (s *QuerySpec) WithQueryFilter(filter *types.QueryFilter) {
	s.queryFilter = filter
}

func (s *QuerySpec) WithCounts() {
	s.includeCounts = true
}

func (s *QuerySpec) Select() (query string, args []any) {
	postgresBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	b := postgresBuilder.Select(s.columns...).From(s.table)
	if s.where != nil {
		b = b.Where(s.where)
	}

	if s.joins != nil {
		for _, join := range s.joins {
			b = b.Join(join)
		}
	}

	if s.groupBys != nil {
		b = b.GroupBy(s.groupBys...)
	}

	if s.orderBy != "" {
		b = b.OrderBy(s.orderBy)
	}

	if s.limit != 0 {
		b = b.Limit(s.limit)
	}

	if s.offset != 0 {
		b = b.Offset(s.offset)
	}

	if s.queryFilter != nil {
		b = applyFilterToQueryBuilder(s.queryFilter, s.table, b)
	}

	query, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}

	return query, args
}

func NewQuery(table string, columns []string) *QuerySpec {
	return &QuerySpec{
		table:   table,
		columns: columns,
	}
}

// applyFilterToQueryBuilder applies the query filter to a query builder.
func applyFilterToQueryBuilder(qf *types.QueryFilter, tableName string, queryBuilder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if qf == nil {
		return queryBuilder
	}

	qf.SetPage(qf.Page)

	if qp := qf.QueryOffset(); qp > 0 {
		queryBuilder = queryBuilder.Offset(uint64(qp))
	}

	if qf.Limit != nil && *qf.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(*qf.Limit))
	} else {
		queryBuilder = queryBuilder.Limit(types.MaxLimit)
	}

	if qf.CreatedAfter != nil {
		queryBuilder = queryBuilder.Where(squirrel.Gt{fmt.Sprintf("%s.%s", tableName, "created_at"): qf.CreatedAfter})
	}

	if qf.CreatedBefore != nil {
		queryBuilder = queryBuilder.Where(squirrel.Lt{fmt.Sprintf("%s.%s", tableName, "created_at"): qf.CreatedBefore})
	}

	if qf.UpdatedAfter != nil {
		queryBuilder = queryBuilder.Where(squirrel.Gt{fmt.Sprintf("%s.%s", tableName, "last_updated_at"): qf.UpdatedAfter})
	}

	if qf.UpdatedBefore != nil {
		queryBuilder = queryBuilder.Where(squirrel.Lt{fmt.Sprintf("%s.%s", tableName, "last_updated_at"): qf.UpdatedBefore})
	}

	return queryBuilder
}

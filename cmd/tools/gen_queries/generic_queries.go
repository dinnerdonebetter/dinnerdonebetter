package main

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func buildListOfNeedingIndexingQuery(table string) string {
	q := goqu.From(table).Select(fmt.Sprintf("%s.id", table)).
		Where(goqu.Ex{
			buildColumnName(table, archivedAt): nil,
		}).
		Where(goqu.Or(
			goqu.Ex{buildColumnName(table, lastIndexedAt): nil},
			goqu.Ex{buildColumnName(table, lastIndexedAt): goqu.Op{exp.LtOp.String(): goqu.L("NOW() - interval '24 hours'")}},
		))

	query, args, err := q.ToSQL()
	if err != nil {
		panic(err)
	}

	if len(args) > 0 {
		panic("args should be empty")
	}

	return query
}

func buildListQuery(table string, columns []string) string {
	anyCols := []any{}
	for _, col := range columns {
		anyCols = append(anyCols, col)
	}

	q := goqu.Select(anyCols...).From(table).
		Where(goqu.Ex{
			buildColumnName(table, archivedAt): nil,
		})

	query, args, err := q.ToSQL()
	if err != nil {
		panic(err)
	}

	if len(args) > 0 {
		panic("args should be empty")
	}

	return query
}

func buildSelectQuery(table string, columns []string) string {
	anyCols := []any{}
	for _, col := range columns {
		anyCols = append(anyCols, col)
	}

	q := goqu.Select(anyCols...).From(table).
		Where(goqu.Ex{
			buildColumnName(table, archivedAt): nil,
		})

	query, _, err := q.ToSQL()
	if err != nil {
		panic(err)
	}

	return query
}

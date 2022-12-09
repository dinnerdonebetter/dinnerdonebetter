package main

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

const (
	id                 = "id"
	dummyValue         = "whatever"
	belongsToHousehold = "belongs_to_household"
	archivedAt         = "archived_at"
)

func buildColumnName(table, column string) string {
	return table + "." + column
}

func mergeSlicesAtIndex[T comparable](a, b []T, index uint) []T {
	return append(a[:index], append(b, a[index:]...)...)
}

func buildJoinStatement(join, on, to string) string {
	return fmt.Sprintf("%s ON %s=%s", join, on, to)
}

func main() {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	fmt.Println(buildGetOneQuery(sqlBuilder))
}

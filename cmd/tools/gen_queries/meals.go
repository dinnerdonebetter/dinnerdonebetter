package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	mealsTable = "meals"
)

/*
var (
	// mealsTableColumns are the columns for the meals table.
	mealsTableColumns = []string{
		buildColumnName(mealsTable, id),
		buildColumnName(mealsTable, "name"),
		buildColumnName(mealsTable, "description"),
		buildColumnName(mealsTable, "min_estimated_portions"),
		buildColumnName(mealsTable, "max_estimated_portions"),
		buildColumnName(mealsTable, "eligible_for_meal_plans"),
		buildColumnName(mealsTable, createdAt),
		buildColumnName(mealsTable, lastUpdatedAt),
		buildColumnName(mealsTable, archivedAt),
		buildColumnName(mealsTable, createdByUser),
	}
)
*/

func buildSelectMealsNeedingIndexingQuery(_ squirrel.StatementBuilderType) string {
	return buildListOfNeedingIndexingQuery(mealsTable)
}

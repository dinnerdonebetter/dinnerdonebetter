package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	recipesTable = "recipes"
)

/*
var (
	recipesTableColumns = []string{
		buildColumnName(recipesTable, id),
		buildColumnName(recipesTable, "name"),
		buildColumnName(recipesTable, "slug"),
		buildColumnName(recipesTable, "source"),
		buildColumnName(recipesTable, "description"),
		buildColumnName(recipesTable, "inspired_by_recipe_id"),
		buildColumnName(recipesTable, "min_estimated_portions"),
		buildColumnName(recipesTable, "max_estimated_portions"),
		buildColumnName(recipesTable, "portion_name"),
		buildColumnName(recipesTable, "plural_portion_name"),
		buildColumnName(recipesTable, "seal_of_approval"),
		buildColumnName(recipesTable, "eligible_for_meals"),
		buildColumnName(recipesTable, "yields_component_type"),
		buildColumnName(recipesTable, createdAt),
		buildColumnName(recipesTable, lastUpdatedAt),
		buildColumnName(recipesTable, archivedAt),
		buildColumnName(recipesTable, "created_by_user"),
	}
)
*/

func buildSelectRecipesNeedingIndexingQuery(_ squirrel.StatementBuilderType) string {
	return buildListOfNeedingIndexingQuery(recipesTable)
}

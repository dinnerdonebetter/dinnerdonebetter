package main

const validIngredientGroupsTableName = "valid_ingredient_groups"

var validIngredientGroupsColumns = []string{
	idColumn,
	"name",
	"slug",
	"description",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

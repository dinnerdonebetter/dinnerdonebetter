package main

const validIngredientGroupsTableName = "valid_ingredient_groups"

var validIngredientGroupsColumns = []string{
	"id",
	"name",
	"slug",
	"description",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

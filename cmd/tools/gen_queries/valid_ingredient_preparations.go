package main

const validIngredientPreparationsTableName = "valid_ingredient_preparations"

var validIngredientPreparationsColumns = []string{
	"id",
	"notes",
	"valid_preparation_id",
	"valid_ingredient_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

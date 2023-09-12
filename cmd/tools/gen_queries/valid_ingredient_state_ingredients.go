package main

const validIngredientStateIngredientsTableName = "valid_ingredient_state_ingredients"

var validIngredientStateIngredientsColumns = []string{
	"id",
	"valid_ingredient",
	"valid_ingredient_state",
	"notes",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

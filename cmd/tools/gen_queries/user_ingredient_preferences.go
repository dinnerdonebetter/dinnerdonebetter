package main

const userIngredientPreferencesTableName = "user_ingredient_preferences"

var userIngredientPreferencesColumns = []string{
	idColumn,
	"ingredient",
	"rating",
	"notes",
	"allergy",
	"belongs_to_user",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

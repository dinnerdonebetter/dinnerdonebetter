package main

const mealComponentsTableName = "meal_components"

var mealComponentsColumns = []string{
	"id",
	"meal_id",
	"recipe_id",
	"meal_component_type",
	"recipe_scale",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

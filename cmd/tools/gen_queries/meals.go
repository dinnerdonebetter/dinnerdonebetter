package main

const mealsTableName = "meals"

var mealsColumns = []string{
	idColumn,
	"name",
	"description",
	"created_by_user",
	"min_estimated_portions",
	"max_estimated_portions",
	"eligible_for_meal_plans",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

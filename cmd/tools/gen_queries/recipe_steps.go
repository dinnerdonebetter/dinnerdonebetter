package main

const recipeStepsTableName = "recipe_steps"

var recipeStepsColumns = []string{
	idColumn,
	"index",
	"preparation_id",
	"minimum_estimated_time_in_seconds",
	"maximum_estimated_time_in_seconds",
	"minimum_temperature_in_celsius",
	"notes",
	"belongs_to_recipe",
	"optional",
	"maximum_temperature_in_celsius",
	"explicit_instructions",
	"condition_expression",
	"start_timer_automatically",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

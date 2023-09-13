package main

const recipeStepVesselsTableName = "recipe_step_vessels"

var recipeStepVesselsColumns = []string{
	idColumn,
	"name",
	"notes",
	"belongs_to_recipe_step",
	"recipe_step_product_id",
	"vessel_predicate",
	"minimum_quantity",
	"maximum_quantity",
	"unavailable_after_step",
	"valid_vessel_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

package main

const recipeStepInstrumentsTableName = "recipe_step_instruments"

var recipeStepInstrumentsColumns = []string{
	idColumn,
	"instrument_id",
	"notes",
	"belongs_to_recipe_step",
	"preference_rank",
	"recipe_step_product_id",
	"name",
	"optional",
	"minimum_quantity",
	"maximum_quantity",
	"option_index",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepInstrumentsQueries() []*Query {
	return []*Query{
		//
	}
}

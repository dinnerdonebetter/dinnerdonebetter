package main

const recipeStepCompletionConditionsTableName = "recipe_step_completion_conditions"

var recipeStepCompletionConditionsColumns = []string{
	idColumn,
	"belongs_to_recipe_step",
	"ingredient_state",
	"notes",
	"optional",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepCompletionConditionQueries() []*Query {
	return []*Query{
		//
	}
}

package main

const recipeStepCompletionConditionIngredientsTableName = "recipe_step_completion_condition_ingredients"

var recipeStepCompletionConditionIngredientsColumns = []string{
	"id",
	"belongs_to_recipe_step_completion_condition",
	"recipe_step_ingredient",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

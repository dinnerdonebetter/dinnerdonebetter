package main

import (
	"github.com/cristalhq/builq"
)

const recipeStepCompletionConditionIngredientsTableName = "recipe_step_completion_condition_ingredients"

var recipeStepCompletionConditionIngredientsColumns = []string{
	idColumn,
	"belongs_to_recipe_step_completion_condition",
	"recipe_step_ingredient",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepCompletionConditionIngredientsQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}

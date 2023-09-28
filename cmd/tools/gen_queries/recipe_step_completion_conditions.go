package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipeStepCompletionConditionsTableName = "recipe_step_completion_conditions"

	ingredientStateColumn = "ingredient_state"
)

var recipeStepCompletionConditionsColumns = []string{
	idColumn,
	"belongs_to_recipe_step",
	ingredientStateColumn,
	notesColumn,
	"optional",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepCompletionConditionQueries() []*Query {
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

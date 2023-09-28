package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipeStepsTableName = "recipe_steps"

	recipeStepIDColumn        = "recipe_step_id"
	belongsToRecipeStepColumn = "belongs_to_recipe_step"
)

var recipeStepsColumns = []string{
	idColumn,
	"index",
	"preparation_id",
	"minimum_estimated_time_in_seconds",
	"maximum_estimated_time_in_seconds",
	"minimum_temperature_in_celsius",
	notesColumn,
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

func buildRecipeStepsQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: "",
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}

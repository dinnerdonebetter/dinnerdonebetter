package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipePrepTaskStepsTableName = "recipe_prep_task_steps"
)

var recipePrepTaskStepsColumns = []string{
	idColumn,
	"satisfies_recipe_step",
	"belongs_to_recipe_step",
	"belongs_to_recipe_prep_task",
}

func buildRecipePrepTaskStepsQueries() []*Query {
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

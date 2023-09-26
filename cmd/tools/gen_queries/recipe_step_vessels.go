package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipeStepVesselsTableName = "recipe_step_vessels"
)

var recipeStepVesselsColumns = []string{
	idColumn,
	nameColumn,
	notesColumn,
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

func buildRecipeStepVesselsQueries() []*Query {
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

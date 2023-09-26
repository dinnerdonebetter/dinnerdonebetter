package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipeStepInstrumentsTableName = "recipe_step_instruments"
)

var recipeStepInstrumentsColumns = []string{
	idColumn,
	"instrument_id",
	notesColumn,
	"belongs_to_recipe_step",
	"preference_rank",
	"recipe_step_product_id",
	nameColumn,
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
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}

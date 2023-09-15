package main

import (
	"github.com/cristalhq/builq"
)

const validIngredientStateIngredientsTableName = "valid_ingredient_state_ingredients"

var validIngredientStateIngredientsColumns = []string{
	idColumn,
	"valid_ingredient",
	"valid_ingredient_state",
	"notes",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientStateIngredientsQueries() []*Query {
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

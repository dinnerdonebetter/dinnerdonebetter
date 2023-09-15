package main

import (
	"github.com/cristalhq/builq"
)

const validIngredientPreparationsTableName = "valid_ingredient_preparations"

var validIngredientPreparationsColumns = []string{
	idColumn,
	"notes",
	"valid_preparation_id",
	"valid_ingredient_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientPreparationsQueries() []*Query {
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

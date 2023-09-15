package main

import (
	"github.com/cristalhq/builq"
)

const validIngredientStatesTableName = "valid_ingredient_states"

var validIngredientStatesColumns = []string{
	idColumn,
	"name",
	"past_tense",
	"slug",
	"description",
	"icon_path",
	"attribute_type",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientStatesQueries() []*Query {
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

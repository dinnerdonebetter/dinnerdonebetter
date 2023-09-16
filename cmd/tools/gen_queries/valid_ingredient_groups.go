package main

import (
	"github.com/cristalhq/builq"
)

const validIngredientGroupsTableName = "valid_ingredient_groups"

var validIngredientGroupsColumns = []string{
	idColumn,
	"name",
	"slug",
	"description",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientGroupsQueries() []*Query {
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

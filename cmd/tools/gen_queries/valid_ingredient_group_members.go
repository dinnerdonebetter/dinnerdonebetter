package main

import (
	"github.com/cristalhq/builq"
)

const validIngredientGroupMembersTableName = "valid_ingredient_group_members"

var validIngredientGroupMembersColumns = []string{
	idColumn,
	"belongs_to_group",
	"valid_ingredient",
	createdAtColumn,
	archivedAtColumn,
}

func buildValidIngredientGroupMembersQueries() []*Query {
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

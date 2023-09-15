package main

import (
	"github.com/cristalhq/builq"
)

const userIngredientPreferencesTableName = "user_ingredient_preferences"

var userIngredientPreferencesColumns = []string{
	idColumn,
	"ingredient",
	"rating",
	"notes",
	"allergy",
	"belongs_to_user",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildUserIngredientPreferencesQueries() []*Query {
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

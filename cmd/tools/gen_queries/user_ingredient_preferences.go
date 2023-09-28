package main

import (
	"github.com/cristalhq/builq"
)

const (
	userIngredientPreferencesTableName = "user_ingredient_preferences"
)

var userIngredientPreferencesColumns = []string{
	idColumn,
	"ingredient",
	"rating",
	notesColumn,
	"allergy",
	belongsToUserColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildUserIngredientPreferencesQueries() []*Query {
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

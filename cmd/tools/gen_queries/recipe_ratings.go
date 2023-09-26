package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipeRatingsTableName = "recipe_ratings"
)

var recipeRatingsColumns = []string{
	idColumn,
	"recipe_id",
	"taste",
	"difficulty",
	"cleanup",
	"instructions",
	"overall",
	notesColumn,
	"by_user",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeRatingsQueries() []*Query {
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

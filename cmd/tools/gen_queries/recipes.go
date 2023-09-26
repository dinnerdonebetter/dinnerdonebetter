package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipesTableName = "recipes"
)

var recipesColumns = []string{
	idColumn,
	nameColumn,
	"source",
	descriptionColumn,
	"inspired_by_recipe_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	"created_by_user",
	"min_estimated_portions",
	"seal_of_approval",
	slugColumn,
	"portion_name",
	"plural_portion_name",
	"max_estimated_portions",
	"eligible_for_meals",
	lastIndexedAtColumn,
	"last_validated_at",
	"yields_component_type",
}

func buildRecipesQueries() []*Query {
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

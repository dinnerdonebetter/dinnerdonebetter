package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipesTableName = "recipes"

	belongsToRecipeColumn = "belongs_to_recipe"
	recipeIDColumn        = "recipe_id"
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
	createdByUserColumn,
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
				Type: "",
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}

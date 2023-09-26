package main

import (
	"github.com/cristalhq/builq"
)

const (
	mealComponentsTableName = "meal_components"
)

var mealComponentsColumns = []string{
	idColumn,
	"meal_id",
	"recipe_id",
	"meal_component_type",
	"recipe_scale",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealComponentsQueries() []*Query {
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

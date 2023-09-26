package main

import (
	"github.com/cristalhq/builq"
)

const (
	mealsTableName = "meals"
)

var mealsColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	"created_by_user",
	"min_estimated_portions",
	"max_estimated_portions",
	"eligible_for_meal_plans",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealsQueries() []*Query {
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

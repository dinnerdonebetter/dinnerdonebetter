package main

import (
	"github.com/cristalhq/builq"
)

const (
	mealPlanOptionsTableName = "meal_plan_options"

	mealPlanOptionIDColumn      = "meal_plan_option_id"
	mealPlanOptionsChosenColumn = "chosen"
)

var mealPlanOptionsColumns = []string{
	idColumn,
	"meal_id",
	notesColumn,
	mealPlanOptionsChosenColumn,
	"tiebroken",
	"assigned_cook",
	"assigned_dishwasher",
	belongsToMealPlanEventColumn,
	"meal_scale",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealPlanOptionsQueries() []*Query {
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

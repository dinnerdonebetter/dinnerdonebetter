package main

import (
	"github.com/cristalhq/builq"
)

const mealPlanEventsTableName = "meal_plan_events"

var mealPlanEventsColumns = []string{
	idColumn,
	"notes",
	"starts_at",
	"ends_at",
	"meal_name",
	"belongs_to_meal_plan",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealPlanEventsQueries() []*Query {
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

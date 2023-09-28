package main

import (
	"github.com/cristalhq/builq"
)

const (
	mealPlansTableName = "meal_plans"

	mealPlanIDColumn     = "meal_plan_id"
	mealPlanStatusColumn = "status"
)

var mealPlansColumns = []string{
	idColumn,
	notesColumn,
	mealPlanStatusColumn,
	"voting_deadline",
	belongsToHouseholdColumn,
	"grocery_list_initialized",
	"tasks_created",
	"election_method",
	createdByUserColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealPlansQueries() []*Query {
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
